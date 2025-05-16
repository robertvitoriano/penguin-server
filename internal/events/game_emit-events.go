package events

type GameEmitEvent string

const (
	SetInitialPlayersPosition GameEmitEvent = "set_initial_players_position"
	UpdatePlayerPosition      GameEmitEvent = "update_player_position"
	MessageReceived           GameEmitEvent = "message_received"
	AudioChunkReceived        GameEmitEvent = "audio_chunk_received"
	WebRTCOfferReceived       GameEmitEvent = "webrtc_offer_received"
	WebRTCCandidateReceived   GameEmitEvent = "webrtc_candidate_received"
	WebRTCAnswerReceived      GameEmitEvent = "webrtc_answer_received"
	PlayerNotFound            GameEmitEvent = "player_not_found"
)
