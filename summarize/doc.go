// Package summarize provides pluggable tool result summarization hooks.
//
// Use summarize.Passthrough for proxy tools before domain summarizers exist.
// Use summarize.Meta to pass context from present hooks to summarizers.
// Limit is the common cross-app field; SinceSeq, RunID, and Extra are optional.
package summarize
