package services

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

type LLMService struct {
	// Add LLM client configurations here
}

func NewLLMService() *LLMService {
	return &LLMService{}
}

type AIResponse struct {
	Response string `json:"response"`
	Error    string `json:"error,omitempty"`
}

type ResumeAnalysis struct {
	Feedback string `json:"feedback"`
	Score    int    `json:"score"`
	Error    string `json:"error,omitempty"`
}

func (s *LLMService) callPythonAI(method string, args map[string]interface{}) ([]byte, error) {
	// Prepare the Python script call
	argsJSON, _ := json.Marshal(args)
	
	cmd := exec.Command("python3", "/home/ubuntu/genai-platform/ai_bridge.py", method, string(argsJSON))
	cmd.Dir = "/home/ubuntu/genai-platform"
	
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("python script error: %v, stderr: %s", err, stderr.String())
	}
	
	return out.Bytes(), nil
}

func (s *LLMService) ProcessPDF(docID int, filePath string) error {
	fmt.Printf("Processing PDF %d at %s\n", docID, filePath)
	
	args := map[string]interface{}{
		"document_id": docID,
		"file_path":   filePath,
	}
	
	_, err := s.callPythonAI("process_document", args)
	if err != nil {
		fmt.Printf("Failed to process PDF %d: %v\n", docID, err)
		return err
	}
	
	fmt.Printf("Successfully processed PDF %d\n", docID)
	return nil
}

func (s *LLMService) GetRelevantContext(documentIDs []int, query string) (string, error) {
	args := map[string]interface{}{
		"query":        query,
		"document_ids": documentIDs,
	}
	
	result, err := s.callPythonAI("search_similar_chunks", args)
	if err != nil {
		return "", err
	}
	
	var chunks []string
	if err := json.Unmarshal(result, &chunks); err != nil {
		return "", err
	}
	
	// Join chunks into context
	context := ""
	for _, chunk := range chunks {
		context += chunk + "\n\n"
	}
	
	return context, nil
}

func (s *LLMService) GenerateResponse(query, context string) (string, error) {
	args := map[string]interface{}{
		"query":   query,
		"context": []string{context},
	}
	
	result, err := s.callPythonAI("generate_chat_response", args)
	if err != nil {
		return "", err
	}
	
	var response AIResponse
	if err := json.Unmarshal(result, &response); err != nil {
		return "", err
	}
	
	if response.Error != "" {
		return "", fmt.Errorf(response.Error)
	}
	
	return response.Response, nil
}

func (s *LLMService) GenerateSQL(naturalQuery string) (string, error) {
	args := map[string]interface{}{
		"natural_query": naturalQuery,
	}
	
	result, err := s.callPythonAI("generate_sql_from_natural_language", args)
	if err != nil {
		return "", err
	}
	
	var response AIResponse
	if err := json.Unmarshal(result, &response); err != nil {
		return "", err
	}
	
	return response.Response, nil
}

func (s *LLMService) ProcessResearchTask(taskID int, query string, db *sql.DB) {
	// Simulate research processing
	time.Sleep(2 * time.Second)
	
	args := map[string]interface{}{
		"research_query": query,
	}
	
	result, err := s.callPythonAI("conduct_research", args)
	if err != nil {
		fmt.Printf("Failed to conduct research for task %d: %v\n", taskID, err)
		return
	}
	
	var response AIResponse
	if err := json.Unmarshal(result, &response); err != nil {
		fmt.Printf("Failed to parse research result for task %d: %v\n", taskID, err)
		return
	}
	
	// Update task with result
	_, err = db.Exec(
		"UPDATE research_tasks SET status = $1, result = $2, completed_at = $3 WHERE id = $4",
		"completed", response.Response, time.Now(), taskID,
	)
	if err != nil {
		fmt.Printf("Failed to update research task %d: %v\n", taskID, err)
	}
}

func (s *LLMService) ProcessResume(analysisID int, resumePath, jobDescription string, db *sql.DB) {
	// Simulate resume processing
	time.Sleep(2 * time.Second)
	
	args := map[string]interface{}{
		"resume_path":     resumePath,
		"job_description": jobDescription,
	}
	
	result, err := s.callPythonAI("analyze_resume", args)
	if err != nil {
		fmt.Printf("Failed to analyze resume for analysis %d: %v\n", analysisID, err)
		return
	}
	
	var analysis ResumeAnalysis
	if err := json.Unmarshal(result, &analysis); err != nil {
		fmt.Printf("Failed to parse resume analysis result for analysis %d: %v\n", analysisID, err)
		return
	}
	
	if analysis.Error != "" {
		fmt.Printf("Resume analysis error for analysis %d: %s\n", analysisID, analysis.Error)
		return
	}
	
	// Update analysis with result
	_, err = db.Exec(
		"UPDATE resume_analyses SET status = $1, feedback = $2, score = $3, completed_at = $4 WHERE id = $5",
		"completed", analysis.Feedback, analysis.Score, time.Now(), analysisID,
	)
	if err != nil {
		fmt.Printf("Failed to update resume analysis %d: %v\n", analysisID, err)
	}
}

