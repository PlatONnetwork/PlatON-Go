package platonstats

import "github.com/PlatONnetwork/PlatON-Go/event"

type SampleEvent struct {
	SampleMsg string
}

type SampleEventProducer struct {
	sampleEventFeed event.Feed
	scope           event.SubscriptionScope
}

func (producer *SampleEventProducer) SubscribeSampleEvent(ch chan<- SampleEvent) event.Subscription {
	return producer.scope.Track(producer.sampleEventFeed.Subscribe(ch))
}
