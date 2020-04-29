package status

import (
	"log"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

}
func Test2FriendsJoining(t *testing.T) {
	incomingCh := make(chan workIn)
	startConsumer(incomingCh)
	defer close(incomingCh)

	ctrl := gomock.NewController(t)

	ClientA := NewMockResponder(ctrl)
	ClientB := NewMockResponder(ctrl)

	//Expectation
	ClientA.EXPECT().IsStateless().Return(false).AnyTimes()
	ClientB.EXPECT().IsStateless().Return(false).AnyTimes()

	ClientA.EXPECT().Reply(&friendResponse{
		UserID: 2,
		Online: false,
	}).Return(nil).Times(1)

	ClientA.EXPECT().Reply(&friendResponse{
		UserID: 2,
		Online: true,
	}).Return(nil).Times(1)

	ClientB.EXPECT().Reply(&friendResponse{
		UserID: 1,
		Online: true,
	}).Return(nil).Times(1)

	incomingCh <- workIn{
		action: Joining,
		payload: &userContext{
			Responder: ClientA,
			ID:        1,
			Friends:   []int{2},
		},
	}

	incomingCh <- workIn{
		action: Joining,
		payload: &userContext{
			Responder: ClientB,
			ID:        2,
			Friends:   []int{1},
		},
	}

	//wait for goroutines to finish - could operate a waitgroup to aid this.
	time.Sleep(20 * time.Millisecond)
	ctrl.Finish()

}

func Test2FriendsJoiningAndLeaving(t *testing.T) {
	incomingCh := make(chan workIn)
	startConsumer(incomingCh)
	defer close(incomingCh)

	ctrl := gomock.NewController(t)

	ClientA := NewMockResponder(ctrl)
	ClientB := NewMockResponder(ctrl)

	//Expectation
	ClientA.EXPECT().IsStateless().Return(false).AnyTimes()
	ClientB.EXPECT().IsStateless().Return(false).AnyTimes()

	ClientA.EXPECT().Reply(&friendResponse{
		UserID: 2,
		Online: false,
	}).Return(nil).Times(1)

	ClientA.EXPECT().Reply(&friendResponse{
		UserID: 2,
		Online: true,
	}).Return(nil).Times(1)

	ClientB.EXPECT().Reply(&friendResponse{
		UserID: 1,
		Online: true,
	}).Return(nil).Times(1)

	ClientB.EXPECT().Reply(&friendResponse{
		UserID: 1,
		Online: false,
	}).Return(nil).Times(1)

	incomingCh <- workIn{
		action: Joining,
		payload: &userContext{
			Responder: ClientA,
			ID:        1,
			Friends:   []int{2},
		},
	}

	incomingCh <- workIn{
		action: Joining,
		payload: &userContext{
			Responder: ClientB,
			ID:        2,
			Friends:   []int{1},
		},
	}

	incomingCh <- workIn{
		action: Leaving,
		payload: &userContext{
			Responder: ClientA,
			ID:        1,
			Friends:   []int{2},
		},
	}

	incomingCh <- workIn{
		action: Leaving,
		payload: &userContext{
			Responder: ClientB,
			ID:        2,
			Friends:   []int{1},
		},
	}

	//wait for goroutines to finish - could operate a waitgroup to aid this.
	time.Sleep(20 * time.Millisecond)
	ctrl.Finish()

}

func TestStatelessTimeout(t *testing.T) {
	incomingCh := make(chan workIn)
	startConsumer(incomingCh)
	defer close(incomingCh)

	ctrl := gomock.NewController(t)

	ClientA := NewMockResponder(ctrl)
	ClientB := NewMockResponder(ctrl)

	//monkey patching
	ttl := OnlineTTL
	OnlineTTL = time.Duration(time.Millisecond * 5)
	defer func() {
		OnlineTTL = ttl
	}()

	//Expectation
	ClientA.EXPECT().IsStateless().Return(true).AnyTimes()
	ClientB.EXPECT().IsStateless().Return(false).AnyTimes()

	ClientA.EXPECT().Reply(&friendResponse{
		UserID: 2,
		Online: false,
	}).Return(nil).Times(1)

	ClientA.EXPECT().Reply(&friendResponse{
		UserID: 2,
		Online: true,
	}).Return(nil).Times(1)

	ClientB.EXPECT().Reply(&friendResponse{
		UserID: 1,
		Online: true,
	}).Return(nil).Times(1)

	//CLIENT A should timeout
	ClientB.EXPECT().Reply(&friendResponse{
		UserID: 1,
		Online: false,
	}).Return(nil).Times(1)

	incomingCh <- workIn{
		action: Joining,
		payload: &userContext{
			Responder: ClientA,
			ID:        1,
			Friends:   []int{2},
		},
	}

	incomingCh <- workIn{
		action: Joining,
		payload: &userContext{
			Responder: ClientB,
			ID:        2,
			Friends:   []int{1},
		},
	}

	//ClientA should time out

	time.Sleep(50 * time.Millisecond)
	ctrl.Finish()
}
