package events

type GameReceiveEvent string

const (
	StartGame      GameReceiveEvent = "start_game"
	PlayerMoved    GameReceiveEvent = "player_moved"
	MessageSent    GameReceiveEvent = "message_sent"
	AudioChuckSent GameReceiveEvent = "audio_chunk_sent"

	WebRTCCandidateSent GameReceiveEvent = "webrtc_candidate_sent"
	WebRTCOfferSent     GameReceiveEvent = "webrtc_offer_sent"
	WebRTCAnswerSent    GameReceiveEvent = "webrtc_answer_sent"
)
