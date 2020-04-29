package status

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func Test2FriendsJoining(t *testing.T) {
	incomingCh := make(chan *requestContext)
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

	incomingCh <- &requestContext{
		statusRequest: &statusRequest{
			UserID:    1,
			FriendIDs: []int{2},
		},
		responder: ClientA,
		action:    Joining,
	}

	incomingCh <- &requestContext{
		statusRequest: &statusRequest{
			UserID:    2,
			FriendIDs: []int{1},
		},
		responder: ClientB,
		action:    Joining,
	}

	//wait for goroutines to finish - could operate a waitgroup to aid this.
	time.Sleep(20 * time.Millisecond)
	ctrl.Finish()

}

func Test2FriendsJoiningAndLeaving(t *testing.T) {
	incomingCh := make(chan *requestContext)
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

	incomingCh <- &requestContext{
		statusRequest: &statusRequest{
			UserID:    1,
			FriendIDs: []int{2},
		},
		responder: ClientA,
		action:    Joining,
	}

	incomingCh <- &requestContext{
		statusRequest: &statusRequest{
			UserID:    2,
			FriendIDs: []int{1},
		},
		responder: ClientB,
		action:    Joining,
	}

	incomingCh <- &requestContext{
		statusRequest: &statusRequest{
			UserID:    1,
			FriendIDs: []int{2},
		},
		responder: ClientA,
		action:    Leaving,
	}

	incomingCh <- &requestContext{
		statusRequest: &statusRequest{
			UserID:    2,
			FriendIDs: []int{1},
		},
		responder: ClientB,
		action:    Leaving,
	}

	//wait for goroutines to finish - could operate a waitgroup to aid this.
	time.Sleep(20 * time.Millisecond)
	ctrl.Finish()

}
