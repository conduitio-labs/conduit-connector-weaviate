package destination

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/conduitio-labs/conduit-connector-weaviate/destination/mock"
	"github.com/conduitio-labs/conduit-connector-weaviate/destination/weaviate"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/google/go-cmp/cmp"
	"github.com/matryer/is"
	"go.uber.org/mock/gomock"
)

type eqMatcher struct {
	want interface{}
}

func newEqMatcher(want interface{}) eqMatcher {
	return eqMatcher{
		want: want,
	}
}

func (eq eqMatcher) Matches(got interface{}) bool {
	return gomock.Eq(eq.want).Matches(got)
}

func (eq eqMatcher) Got(got interface{}) string {
	diff := cmp.Diff(
		got,
		eq.want,
	)
	return fmt.Sprintf(
		"%v (%T)\nDiff (-got +want):\n%s",
		got,
		got,
		strings.TrimSpace(diff),
	)
}

func (eq eqMatcher) String() string {
	return fmt.Sprintf("%v (%T)\n", eq.want, eq.want)
}

func TestDestination_Teardown_NoOpen(t *testing.T) {
	is := is.New(t)
	underTest := New()
	err := underTest.Teardown(context.Background())
	is.NoErr(err)
}

func TestDestination_Open_OpensClient(t *testing.T) {
	ctx := context.Background()
	cfg := map[string]string{
		"endpoint":           "test-endpoint",
		"scheme":             "test-scheme",
		"apiKey":             "test-api-key",
		"class":              "test-class",
		"moduleAPIKey.name":  "X-OpenAI-Api-Key",
		"moduleAPIKey.value": "test-OpenAI-Api-Key",
		"generateUUID":       "true",
	}

	// setupTest is doing the basic checks
	_, _ = setupTest(t, ctx, cfg)
}
func TestDestination_SingleWrite(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	cfg := map[string]string{
		"endpoint":           "test-endpoint",
		"scheme":             "test-scheme",
		"apiKey":             "test-api-key",
		"class":              "test-class",
		"moduleAPIKey.name":  "X-OpenAI-Api-Key",
		"moduleAPIKey.value": "test-OpenAI-Api-Key",
		"generateUUID":       "false",
	}

	testCases := []struct {
		name   string
		record sdk.Record
		want   *weaviate.Object
	}{
		{
			name: "raw payload, raw key, use class from config",
			record: sdk.Util.Source.NewRecordCreate(
				sdk.Position("test-position"),
				map[string]string{},
				sdk.RawData("f9a510b3-5865-40e4-9fe8-e7fbab25b8bc"),
				sdk.RawData(`{
					"product_name": "computer",
					"price":        1000,
					"labels":       ["laptop", "navy-blue"],
					"used":			true,
					"specs": 		{
						"cpu": 		"3GHz",
						"memory":	"16Gb"
					}
				}`),
			),
			want: &weaviate.Object{
				ID:    "f9a510b3-5865-40e4-9fe8-e7fbab25b8bc",
				Class: cfg["class"],
				Properties: map[string]interface{}{
					"product_name": "computer",
					"price":        float64(1000),
					"labels":       []any{"laptop", "navy-blue"},
					"used":         true,
					"specs": map[string]any{
						"cpu":    "3GHz",
						"memory": "16Gb",
					},
				},
			},
		},
		{
			name: "structured payload, raw key, use class from config",
			record: sdk.Util.Source.NewRecordCreate(
				sdk.Position("test-position"),
				map[string]string{},
				sdk.RawData("f9a510b3-5865-40e4-9fe8-e7fbab25b8bc"),
				sdk.StructuredData{
					"product_name": "computer",
					"price":        1000,
					"labels":       []string{"laptop", "navy-blue"},
					"used":         true,
					"specs": map[string]any{
						"cpu":    "3GHz",
						"memory": "16Gb",
					},
				},
			),
			want: &weaviate.Object{
				ID:    "f9a510b3-5865-40e4-9fe8-e7fbab25b8bc",
				Class: cfg["class"],
				Properties: map[string]interface{}{
					"product_name": "computer",
					"price":        float64(1000),
					"labels":       []any{"laptop", "navy-blue"},
					"used":         true,
					"specs": map[string]any{
						"cpu":    "3GHz",
						"memory": "16Gb",
					},
				},
			},
		},
		{
			name: "structured payload, raw key, use class from metadata",
			record: sdk.Util.Source.NewRecordCreate(
				sdk.Position("test-position"),
				map[string]string{
					metadataClass: "top-secret-class",
				},
				sdk.RawData("f9a510b3-5865-40e4-9fe8-e7fbab25b8bc"),
				sdk.StructuredData{
					"product_name": "computer",
					"price":        1000,
					"labels":       []string{"laptop", "navy-blue"},
					"used":         true,
					"specs": map[string]any{
						"cpu":    "3GHz",
						"memory": "16Gb",
					},
				},
			),
			want: &weaviate.Object{
				ID:    "f9a510b3-5865-40e4-9fe8-e7fbab25b8bc",
				Class: "top-secret-class",
				Properties: map[string]interface{}{
					"product_name": "computer",
					"price":        float64(1000),
					"labels":       []any{"laptop", "navy-blue"},
					"used":         true,
					"specs": map[string]any{
						"cpu":    "3GHz",
						"memory": "16Gb",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			underTest, wClient := setupTest(t, ctx, cfg)
			wClient.EXPECT().Insert(ctx, newEqMatcher(tc.want))

			n, err := underTest.Write(ctx, []sdk.Record{tc.record})
			is.NoErr(err)
			is.Equal(1, n)
		})
	}
}

func setupTest(t *testing.T, ctx context.Context, cfg map[string]string) (sdk.Destination, *mock.WeaviateClient) {
	is := is.New(t)
	ctrl := gomock.NewController(t)
	client := mock.NewWeaviateClient(ctrl)
	client.EXPECT().
		Open(gomock.Eq(weaviate.Config{
			APIKey:   cfg["apiKey"],
			Endpoint: cfg["endpoint"],
			Scheme:   cfg["scheme"],
			Headers: map[string]string{
				"X-OpenAI-Api-Key": "test-OpenAI-Api-Key",
			},
		}))

	underTest := NewWithClient(client)
	err := underTest.Configure(ctx, cfg)
	is.NoErr(err)

	err = underTest.Open(ctx)
	is.NoErr(err)

	return underTest, client
}
