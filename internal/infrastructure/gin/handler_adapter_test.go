package gin_test

import (
	"bytes"
	"context"
	"errors"
	ginutils "gochat/internal/infrastructure/gin"
	"gochat/internal/infrastructure/gin/mocks"
	"gochat/internal/shared/api"
	myErrors "gochat/internal/shared/errors"
	httptestutil "gochat/pkg/httptest"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type testRequest struct {
	Name string `json:"name"`
}

func (r *testRequest) Bind(c *gin.Context) error {
	return c.BindJSON(r)
}

type testInput struct{ Name string }

func (i *testInput) Validate() error {
	if i.Name == "invalid" {
		return errors.New("invalid input")
	}
	return nil
}

type testOutput struct {
	Greeting string `json:"greeting"`
}

type fakeUseCase struct {
	exec func(ctx context.Context, in *testInput) (*testOutput, error)
}

func (f fakeUseCase) Execute(ctx context.Context, in *testInput) (*testOutput, error) {
	return f.exec(ctx, in)
}

func Test_AdaptUseCaseToHandler_Flows(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		body            string
		expectValidator func(validator *mocks.MockValidator)
		useCase         fakeUseCase
		expectResponse  *api.Response
	}{
		{
			name: "success",
			body: `{"name":"Alice"}`,
			expectValidator: func(validator *mocks.MockValidator) {
				validator.EXPECT().Validate(gomock.Any(), gomock.Any()).Return("", nil)
			},
			useCase: fakeUseCase{exec: func(_ context.Context, in *testInput) (*testOutput, error) {
				assert.Equal(t, "Alice", in.Name)
				return &testOutput{Greeting: "hi Alice"}, nil
			}},
			expectResponse: api.NewResponseWithData(&testOutput{Greeting: "hi Alice"}),
		},
		{
			name:            "bind error -> invalid param",
			body:            `{"name":`, // 错误 JSON
			expectValidator: nil,
			useCase:         fakeUseCase{exec: func(_ context.Context, in *testInput) (*testOutput, error) { return nil, nil }},
			expectResponse:  api.ResponseInvalidParam,
		},
		{
			name: "validator error -> server error",
			body: `{"name":"Alice"}`,
			expectValidator: func(validator *mocks.MockValidator) {
				validator.EXPECT().Validate(gomock.Any(), gomock.Any()).Return("", errors.New("boom"))
			},
			useCase:        fakeUseCase{exec: func(_ context.Context, in *testInput) (*testOutput, error) { return nil, nil }},
			expectResponse: api.ResponseServerError,
		},
		{
			name: "validator message -> invalid param with message",
			body: `{"name":""}`,
			expectValidator: func(validator *mocks.MockValidator) {
				validator.EXPECT().Validate(gomock.Any(), gomock.Any()).Return("field name is required", nil)
			},
			useCase:        fakeUseCase{exec: func(_ context.Context, in *testInput) (*testOutput, error) { return nil, nil }},
			expectResponse: api.NewResponseWithMessage(api.CodeInvalidParam, "field name is required"),
		},
		{
			name: "input.Validate error -> invalid param",
			body: `{"name":"invalid"}`,
			expectValidator: func(validator *mocks.MockValidator) {
				validator.EXPECT().Validate(gomock.Any(), gomock.Any()).Return("", nil)
			},
			useCase: fakeUseCase{exec: func(_ context.Context, in *testInput) (*testOutput, error) {
				t.Fatalf("use case should not be called")
				return nil, nil
			}},
			expectResponse: api.ResponseInvalidParam,
		},
		{
			name: "use case business error -> business code & message",
			body: `{"name":"Alice"}`,
			expectValidator: func(validator *mocks.MockValidator) {
				validator.EXPECT().Validate(gomock.Any(), gomock.Any()).Return("", nil)
			},
			useCase: fakeUseCase{exec: func(_ context.Context, in *testInput) (*testOutput, error) {
				return nil, myErrors.NewBusiness("something wrong")
			}},
			expectResponse: api.NewResponseWithMessage(api.CodeBusinessError, "something wrong"),
		},
		{
			name: "use case system error -> server error",
			body: `{"name":"Alice"}`,
			expectValidator: func(validator *mocks.MockValidator) {
				validator.EXPECT().Validate(gomock.Any(), gomock.Any()).Return("", nil)
			},
			useCase: fakeUseCase{exec: func(_ context.Context, in *testInput) (*testOutput, error) {
				return nil, errors.New("io")
			}},
			expectResponse: api.ResponseServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			router := httptestutil.NewTestRouter(t)

			h := ginutils.AdaptUseCaseToHandler[
				testRequest,
				*testRequest,
				*testInput,
				*testOutput,
			](
				fakeUseCase{exec: func(ctx context.Context, in *testInput) (*testOutput, error) {
					return tt.useCase.exec(ctx, in)
				}},
				func() ginutils.Validator {
					validator := mocks.NewMockValidator(ctrl)
					if tt.expectValidator != nil {
						tt.expectValidator(validator)
					}
					return validator
				}(),
				func(req *testRequest) *testInput { return &testInput{Name: req.Name} },
				func(c *gin.Context, out *testOutput) {
					ginutils.ResponseSuccessWithData(c, out)
				},
			)

			router.POST("/test", h)

			httptestutil.ExecuteRouteTest(t, router, &httptestutil.TestCase{
				Path:     "/test",
				Method:   http.MethodPost,
				Body:     bytes.NewBufferString(tt.body),
				Response: tt.expectResponse,
			})
		})
	}
}
