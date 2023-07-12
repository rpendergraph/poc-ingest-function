package moveit

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fn/elasticsearch"
	"fn/types"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Event struct {
	AgentGUID         string    `json:"agentGuid"`
	DestHost          string    `json:"destHost"`
	DestHostType      int       `json:"destHostType"`
	DestID            int       `json:"destId"`
	DestLocation      string    `json:"destLocation"`
	Direction         string    `json:"direction"`
	Duration          int       `json:"duration"`
	EntityID          int       `json:"entityId"`
	Identity          string    `json:"identity"`
	IdentityID        int       `json:"identityId"`
	InitialID         int       `json:"initialId"`
	ObjectDestBytes   int       `json:"objectDestBytes"`
	ObjectDestName    string    `json:"objectDestName"`
	ObjectSourceBytes int       `json:"objectSourceBytes"`
	ObjectSourceName  string    `json:"objectSourceName"`
	Operation         string    `json:"operation"`
	OperationType     string    `json:"operationType"`
	OrgID             string    `json:"orgId"`
	OrgName           string    `json:"orgName"`
	ServerName        string    `json:"serverName"`
	ServerSubType     string    `json:"serverSubType"`
	ServerType        string    `json:"serverType"`
	SourceHost        string    `json:"sourceHost"`
	SourceHostType    string    `json:"sourceHostType"`
	SourceID          string    `json:"sourceId"`
	SourceLocation    string    `json:"sourceLocation"`
	Status            string    `json:"status"`
	StatusMsg         string    `json:"statusMsg"`
	Timestamp         time.Time `json:"timestamp"`
	WorkflowID        string    `json:"workflowId"`
}

type IndexRecord struct {
	Governance types.GovernedEvent `json:"governance"`
	MoveIt     Event               `json:"moveit"`
	Timestamp  string              `json:"@timestamp"`
}

type handler struct {
	writer elasticsearch.IndexWriter[IndexRecord]
}

func NewMoveItHandler(indexWriter elasticsearch.IndexWriter[IndexRecord]) (http.Handler, error) {
	return handler{
		writer: indexWriter,
	}, nil
}

func (r handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	events, err := parseMoveItBody(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	indexDocuments := r.createIndexedDocuments(events)

	err = r.writer.WriteToIndex(req.Context(), indexDocuments)
	if err != nil {
		log.Printf("Error writing to index: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (r handler) createIndexedDocuments(events []Event) []IndexRecord {
	indexDocuments := make([]IndexRecord, len(events))

	for i, e := range events {
		governance := r.createDataGovernanceRecord(e)

		indexDocuments[i] = IndexRecord{
			Governance: governance,
			MoveIt:     e,
			Timestamp:  e.Timestamp.String(),
		}
	}

	return indexDocuments
}

func (r handler) createDataGovernanceRecord(e Event) types.GovernedEvent {
	return types.GovernedEvent{
		ID:             e.AgentGUID,
		ContainsPII:    false,
		ContainsPHI:    true,
		LineOfBusiness: "FakeLOB",
		Source:         e.ObjectSourceName,
		Destination:    e.ObjectDestName,
		SizeBytes:      uint(e.ObjectDestBytes),
		Tags:           []string{"Tag1, Tag2"},
	}
}

func parseMoveItBody(body []byte) ([]Event, error) {
	bodyReader := bytes.NewReader(body)
	scanner := bufio.NewScanner(bodyReader)
	scanner.Split(bufio.ScanLines)
	var out []Event
	for scanner.Scan() {
		var event Event
		line := strings.ReplaceAll(strings.Trim(scanner.Text(), `"`), `""`, `"`)
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return nil, err
		}
		out = append(out, event)
	}
	return out, nil
}
