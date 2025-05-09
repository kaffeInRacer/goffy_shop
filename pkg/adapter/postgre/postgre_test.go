package postgre_test

import (
	"context"
	"testing"

	"github.com/kaffein/goffy/pkg/adapter/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestMockAdapter_StartStop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapter := mock.NewMockAdapter(ctrl)
	mockAdapter.EXPECT().Start(gomock.Any()).Return(nil).Times(1)
	mockAdapter.EXPECT().Stop(gomock.Any()).Return(nil).Times(1)

	err := mockAdapter.Start(context.Background())
	assert.NoError(t, err)

	err = mockAdapter.Stop(context.Background())
	assert.NoError(t, err)
}
