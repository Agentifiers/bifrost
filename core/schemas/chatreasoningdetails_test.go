package schemas

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestChatReasoningDetailsFormatSurvivesResponseRoundTrip(t *testing.T) {
	const details = `[{"type":"reasoning.text","text":"step by step","signature":"sig-abc","format":"anthropic-claude-v1","index":0},{"type":"reasoning.encrypted","data":"enc-xyz","id":"rd_1","format":"anthropic-claude-v1","index":1}]`

	tests := []struct {
		name     string
		response string
		path     string
	}{
		{
			name:     "non-stream message",
			response: `{"id":"gen-1","object":"chat.completion","created":1,"model":"anthropic/claude-sonnet-4.5","choices":[{"index":0,"finish_reason":"stop","message":{"role":"assistant","content":"It is noon.","reasoning":"step by step","reasoning_details":` + details + `}}]}`,
			path:     "choices.0.message.reasoning_details",
		},
		{
			name:     "stream delta",
			response: `{"id":"gen-1","object":"chat.completion.chunk","created":1,"model":"anthropic/claude-sonnet-4.5","choices":[{"index":0,"finish_reason":null,"delta":{"role":"assistant","reasoning":"step by step","reasoning_details":` + details + `}}]}`,
			path:     "choices.0.delta.reasoning_details",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var response BifrostChatResponse
			require.NoError(t, Unmarshal([]byte(tt.response), &response))

			encoded, err := Marshal(response)
			require.NoError(t, err)

			require.JSONEq(t, details, gjson.GetBytes(encoded, tt.path).Raw)
		})
	}
}
