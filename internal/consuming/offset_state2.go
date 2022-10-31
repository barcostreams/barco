package consuming

import (
	"fmt"
	"sync"

	"github.com/barcostreams/barco/internal/conf"
	"github.com/barcostreams/barco/internal/data"
	"github.com/barcostreams/barco/internal/discovery"
	"github.com/barcostreams/barco/internal/interbroker"
	"github.com/barcostreams/barco/internal/localdb"
	. "github.com/barcostreams/barco/internal/types"
	"github.com/rs/zerolog/log"
)

func newDefaultOffsetState2(
	localDb localdb.Client,
	discoverer discovery.TopologyGetter,
	gossiper interbroker.Gossiper,
	config conf.ConsumerConfig,
) OffsetState {
	state := &defaultOffsetState2{
		offsetMap:  make(map[OffsetStoreKey]*Offset),
		commitChan: make(chan *OffsetStoreKeyValue, 64),
		localDb:    localDb,
		gossiper:   gossiper,
		discoverer: discoverer,
		config:     config,
	}
	go state.processCommit()
	return state
}

// Stores offsets by range, this
type defaultOffsetState2 struct {
	offsetMap  map[OffsetStoreKey]*Offset
	mu         sync.RWMutex              // We need synchronization for doing CAS operation per key/value
	commitChan chan *OffsetStoreKeyValue // We need to commit offset in order
	localDb    localdb.Client
	gossiper   interbroker.Gossiper
	discoverer discovery.TopologyGetter
	config     conf.ConsumerConfig
}

func (s *defaultOffsetState2) Init() error {
	// Load local offsets into memory
	values, err := s.localDb.Offsets()
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, kv := range values {
		s.offsetMap[kv.Key] = &kv.Value
	}

	return nil
}

func (s *defaultOffsetState2) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return fmt.Sprint(s.offsetMap)
}

func (s *defaultOffsetState2) Get(group string, topic string, token Token, rangeIndex RangeIndex, clusterSize int) (offset *Offset, rangesMatch bool) {
	key := OffsetStoreKey{Group: group, Topic: topic}
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.offsetMap[key], true
}

func (s *defaultOffsetState2) Set(
	group string,
	topic string,
	value Offset,
	commit OffsetCommitType,
) {
	key := OffsetStoreKey{Group: group, Topic: topic}
	// We could use segment logs in the future
	s.mu.Lock()
	defer s.mu.Unlock()
	existingValue := s.offsetMap[key]

	if s.isOldValue(existingValue, &value) {
		// We have a newer value in our map, don't override
		return
	}

	s.offsetMap[key] = &value

	if commit != OffsetCommitNone {
		// Store commits locally in order but don't await for it to complete
		kv := &OffsetStoreKeyValue{Key: key, Value: value}
		s.commitChan <- kv

		if commit == OffsetCommitAll {
			// Send to followers in the background with no order guarantees
			// The local OffsetState of the follower will verify for new values
			go s.sendToFollowers(kv)
		}
	}
}

func (s *defaultOffsetState2) processCommit() {
	for kv := range s.commitChan {
		if err := s.localDb.SaveOffset(kv); err != nil {
			log.Err(err).Interface("offset", *kv).Msgf("Offset could not be stored in the local db")
		}
	}
}

func (s *defaultOffsetState2) isOldValue(existing *Offset, newValue *Offset) bool {
	// TODO: USE TIMESTAMP
	if existing == nil {
		// There's no previous value
		return false
	}

	if existing.Source.Id.Start == newValue.Source.Id.Start {
		// Same tokens (most common case)
		if existing.Source.Id.Version < newValue.Source.Id.Version {
			// The source of the previous value is old
			return false
		}

		if existing.Source.Id.Version > newValue.Source.Id.Version {
			// The new value's source is old
			return true
		}
	}

	if existing.Version < newValue.Version {
		return false
	}
	if existing.Version == newValue.Version && existing.Offset <= newValue.Offset {
		return false
	}

	return true
}

func (s *defaultOffsetState2) sendToFollowers(kv *OffsetStoreKeyValue) {
	id := GenId{Start: kv.Value.Token, Version: kv.Value.Version}
	gen := s.discoverer.GenerationInfo(id)
	if gen == nil {
		log.Error().
			Interface("id", id).
			Msgf("Generation could not be retrieved when saving offset")
		return
	}

	topology := s.discoverer.Topology()

	for _, follower := range gen.Followers {
		if follower == topology.MyOrdinal() {
			continue
		}

		ordinal := follower
		go func() {
			err := s.gossiper.SendCommittedOffset(ordinal, kv)
			if err != nil {
				if topology.HasBroker(ordinal) {
					log.Err(err).Msgf("Offset could not be sent to follower B%d", ordinal)
				}
			} else {
				log.Debug().Msgf(
					"Offset sent to follower B%d for group %s topic '%s' %d/%d",
					ordinal, kv.Key.Group, kv.Key.Topic, kv.Value.Token, kv.Value.Index)
			}
		}()
	}
}

func (s *defaultOffsetState2) ProducerOffsetLocal(topic *TopicDataId) (int64, error) {
	return data.ReadProducerOffset(topic, s.config)
}

func (s *defaultOffsetState2) MinOffset(topic string, token Token, index RangeIndex) *Offset {
	// TODO: CHECK WHETHER IT CAN BE REMOVED
	return nil
}
