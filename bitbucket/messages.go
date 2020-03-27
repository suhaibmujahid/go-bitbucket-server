// This file provides functions for validating payloads from Bitbucket Server Webhooks.
// Doc: https://confluence.atlassian.com/bitbucketserver070/managing-webhooks-in-bitbucket-server-996644364.html

package bitbucket

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	// sha1Prefix is the prefix used by Bitbucket Server before the HMAC hexdigest.
	sha1Prefix = "sha1"
	// sha256Prefix and sha512Prefix are provided for future compatibility.
	sha256Prefix = "sha256"
	sha512Prefix = "sha512"
	// signatureHeader is the Bitbucket Server header key used to pass the HMAC hexdigest.
	signatureHeader = "X-Hub-Signature"
	// eventKeyHeader is the Bitbucket Server header key used to pass the event type.
	eventKeyHeader = "X-Event-Key"
	// requestIDHeader is the Bitbucket Server header key used to pass the unique ID for the webhook event.
	requestIDHeader = "X-Request-Id"
	// payloadFormParam is the name of the form parameter that the JSON payload
	// will be in if a webhook has its content type set to application/x-www-form-urlencoded.
	payloadFormParam = "payload"
)

// genMAC generates the HMAC signature for a message provided the secret key
// and hashFunc.
func genMAC(message, key []byte, hashFunc func() hash.Hash) []byte {
	mac := hmac.New(hashFunc, key)
	mac.Write(message)
	return mac.Sum(nil)
}

// checkMAC reports whether messageMAC is a valid HMAC tag for message.
func checkMAC(message, messageMAC, key []byte, hashFunc func() hash.Hash) bool {
	expectedMAC := genMAC(message, key, hashFunc)
	return hmac.Equal(messageMAC, expectedMAC)
}

// messageMAC returns the hex-decoded HMAC tag from the signature and its
// corresponding hash function.
func messageMAC(signature string) ([]byte, func() hash.Hash, error) {
	if signature == "" {
		return nil, nil, errors.New("missing signature")
	}
	sigParts := strings.SplitN(signature, "=", 2)
	if len(sigParts) != 2 {
		return nil, nil, fmt.Errorf("error parsing signature %q", signature)
	}

	var hashFunc func() hash.Hash
	switch sigParts[0] {
	case sha1Prefix:
		hashFunc = sha1.New
	case sha256Prefix:
		hashFunc = sha256.New
	case sha512Prefix:
		hashFunc = sha512.New
	default:
		return nil, nil, fmt.Errorf("unknown hash type prefix: %q", sigParts[0])
	}

	buf, err := hex.DecodeString(sigParts[1])
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding signature %q: %v", signature, err)
	}
	return buf, hashFunc, nil
}

// ValidatePayload validates an incoming Bitbucket Server Webhook event request
// and returns the (JSON) payload.
// The Content-Type header of the payload can be "application/json" or "application/x-www-form-urlencoded".
// If the Content-Type is neither then an error is returned.
// secretToken is the Bitbucket Server Webhook secret token.
// If your webhook does not contain a secret token, you can pass nil or an empty slice.
// This is intended for local development purposes only and all webhooks should ideally set up a secret token.
//
// Example usage:
//
//     func (s *BitbucketEventMonitor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//       payload, err := bitbucket.ValidatePayload(r, s.webhookSecretKey)
//       if err != nil { ... }
//       // Process payload...
//     }
//
func ValidatePayload(r *http.Request, secretToken []byte) (payload []byte, err error) {
	var body []byte // Raw body that Bitbucket Server uses to calculate the signature.

	switch ct := r.Header.Get("Content-Type"); ct {
	case "application/json":
		var err error
		if body, err = ioutil.ReadAll(r.Body); err != nil {
			return nil, err
		}

		// If the content type is application/json,
		// the JSON payload is just the original body.
		payload = body

	case "application/x-www-form-urlencoded":
		var err error
		if body, err = ioutil.ReadAll(r.Body); err != nil {
			return nil, err
		}

		// If the content type is application/x-www-form-urlencoded,
		// the JSON payload will be under the "payload" form param.
		form, err := url.ParseQuery(string(body))
		if err != nil {
			return nil, err
		}
		payload = []byte(form.Get(payloadFormParam))

	default:
		return nil, fmt.Errorf("webhook request has unsupported Content-Type %q", ct)
	}

	// Only validate the signature if a secret token exists. This is intended for
	// local development only and all webhooks should ideally set up a secret token.
	if len(secretToken) > 0 {
		sig := r.Header.Get(signatureHeader)
		if err := ValidateSignature(sig, body, secretToken); err != nil {
			return nil, err
		}
	}

	return payload, nil
}

// ValidateSignature validates the signature for the given payload.
// signature is the Bitbucket Server hash signature delivered in the X-Hub-Signature header.
// payload is the JSON payload sent by Bitbucket Server Webhooks.
// secretToken is the Bitbucket Server Webhook secret token.
//
// Doc: https://confluence.atlassian.com/bitbucketserver070/managing-webhooks-in-bitbucket-server-996644364.html#ManagingwebhooksinBitbucketServer-webhooksecretsWebhooksecrets
func ValidateSignature(signature string, payload, secretToken []byte) error {
	messageMAC, hashFunc, err := messageMAC(signature)
	if err != nil {
		return err
	}
	if !checkMAC(payload, messageMAC, secretToken, hashFunc) {
		return errors.New("payload signature check failed")
	}
	return nil
}

// WebHookType returns the event key of webhook request r.
//
// Doc: https://confluence.atlassian.com/bitbucketserver070/event-payload-996644369.html#Eventpayload-HTTPheaders
func WebHookType(r *http.Request) string {
	return r.Header.Get(eventKeyHeader)
}

// RequestID returns the unique UUID for each webhook request r.
//
// Doc: https://confluence.atlassian.com/bitbucketserver070/event-payload-996644369.html#Eventpayload-HTTPheaders
func RequestID(r *http.Request) string {
	return r.Header.Get(requestIDHeader)
}

// ParseWebHook parses the event payload. An error will be returned for unrecognized
// event types.
//
// Example usage:
//
//     func (s *BitbucketEventMonitor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//       payload, err := bitbucket.ValidatePayload(r, s.webhookSecretKey)
//       if err != nil { ... }
//       event, err := bitbucket.ParseWebHook(bitbucket.WebHookType(r), payload)
//       if err != nil { ... }
//       switch event := event.(type) {
//       case *bitbucket.RepositoryForkedEvent:
//           processRepositoryForkedEvent(event)
//       case *bitbucket.PullRequestModifiedEvent:
//           processPullRequestModifiedEvent(event)
//       ...
//       }
//     }
//
func ParseWebHook(eventKey string, payload []byte) (interface{}, error) {
	var event interface{}

	switch eventKey {
	case EventKeyRepositoryPush:
		event = &PushEvent{}
	case EventKeyRepositoryModified:
		event = &RepositoryModifiedEvent{}
	case EventKeyRepositoryForked:
		event = &RepositoryForkedEvent{}
	case EventKeyPullRequestOpened:
		event = &PullRequestOpenedEvent{}
	case EventKeyPullRequestReviewersUpdated:
		event = &PullRequestReviewerEvent{}
	case EventKeyPullRequestModified:
		event = &PullRequestModifiedEvent{}
	case EventKeyPullRequestBranchUpdated:
		event = &PullRequestBranchUpdatedEvent{}
	case EventKeyPullRequestApproved:
		event = &PullRequestApprovedEvent{}
	case EventKeyPullRequestUnapproved:
		event = &PullRequestUnapprovedEvent{}
	case EventKeyPullRequestNeedsWork:
		event = &PullRequestNeedsWorkEvent{}
	case EventKeyPullRequestMerged:
		event = &PullRequestMergedEvent{}
	case EventKeyPullRequestDeclined:
		event = &PullRequestDeclinedEvent{}
	case EventKeyPullRequestDeleted:
		event = &PullRequestDeletedEvent{}

	default:
		return nil, fmt.Errorf("unknown X-Event-Key in message: %v", eventKey)
	}

	err := json.Unmarshal(payload, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}
