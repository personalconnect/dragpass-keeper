package main

import (
	"io"
	"log"
	"os"

	"github.com/personalconnect/dragpass-keeper/internal/keystore"
)

// keystore에 저장되는 항목들:
// - 서버 공개키 (init시 저장)
// - 디바이스키
// - 세션코드
// - Keeper 비공개키
// - Keeper 공개키

// API Actions:
// (ping) 헬스 체크
// (savedevicekey) 디바이스키 저장 요청
// (deletedevicekey) 디바이스키 삭제 요청
// (getdevicekey) 디바이스키 조회 요청
// (generatekeypair) 키페어 생성 요청 [Internal: 세션 코드 삭제, 기존 키페어 삭제, 새 키페어 저장]
// (getsessioncode) 세션코드 조회 요청
// (getpublickey) Keeper 공개키 조회 요청

// 회원가입:
// (signalias) Alias를 전달 -> Alias에 Helper 비공개키로 Signature 생성 -> Signature, Helper 공개키 반환
// (savesessioncode) 암호화된 세션 코드, Signature -> 서버 공개키로 Signature 검증, Helper 비공개키로 복호화, 세션 코드 저장 -> 세션 코드 반환

// 로그인:
// (signaliaswithtimestamp) Alias를 전달 -> Alias + Timestamp에 Helper 비공개키로 Signature 생성 (서명) -> Signature, Timestamp 반환
// (signchallengetoken) Signature, ChallengeToken 전달 -> 서버 공개키로 Signature 검증, Helper 비공개키로 챌린지 토큰에 서명 -> Signature 반환

// 다른 기기에 로그인:
// (generatekeypair) Signature, ChallengeToken 전달 -> 서버 공개키로 Signature 검증, 키페어 생성 -> Public Key 반환
// (getpublickey) Keeper 공개키 조회
// (savesessioncode) 암호화된 세션 코드 저장

func init() {
	if err := keystore.EnsureServerPublicKey(); err != nil {
		log.Fatalf("Critical: Failed to ensure server public key: %v", err)
	}
}

func main() {
	// Stdout is sent to the Chrome extension, so we log to Stderr
	log.SetOutput(os.Stderr)

	// For debugging, log to a file
	// logFile, _ := os.OpenFile("/tmp/keeper.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// log.SetOutput(logFile)

	if err := keystore.LoadBinaryInfo(); err != nil {
		log.Printf("Warning: Failed to calculate binary info: %v", err)
	}

	if err := keystore.EnsureServerPublicKey(); err != nil {
		log.Fatalf("Critical: Failed to ensure server public key: %v", err)
	}

	log.Println("DragPass extension helper started")
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Critical Panic Recovered: %v", r)
		}
	}()

	msgr := keystore.NewMessenger(os.Stdin, os.Stdout)
	for {

		// Read raw message bytes
		msg, err := msgr.ReadMessage()
		if err != nil {
			if err == io.EOF {
				log.Println("Chrome extension closed the connection")
				break
			}
			log.Printf("Failed to read message: %v", err)

			errorResponse := keystore.BaseResponse{
				Success: false,
				Error:   "Native host read error: " + err.Error(),
			}

			if sendErr := msgr.SendResponse(errorResponse); sendErr != nil {
				log.Printf("Failed to send error response: %v", sendErr)
				break
			}
			continue
		}

		// Handle the request
		resp := keystore.HandleRequest(msg)

		// Send response
		if err := msgr.SendResponse(resp); err != nil {
			log.Printf("Failed to send response: %v", err)
			break
		}
	}
}
