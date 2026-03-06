// Package agent provides core agent scheduling logic.
//
// This package contains the Agent struct which orchestrates multi-turn conversations by coordinating between the LLM
// client, session store, and user interface. It handles the complete chat workflow including loading sessions,
// streaming responses from the LLM, and persisting conversation history.
package agent
