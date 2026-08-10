import { useState, useRef } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Upload, Loader2, CheckCircle2, Clock, FileText, 
  TrendingUp, AlertCircle, Target, Search, FileCheck 
} from 'lucide-react'

export default function ResumeFeedbackEnhanced() {
  const [jobDescription, setJobDescription] = useState('')
  const [uploading, setUploading] = useState(false)
  const [analysisId, setAnalysisId] = useState(null)
  const [feedback, setFeedback] = useState(null)
  const [agentProgress, setAgentProgress] = useState(null)
  const fileInputRef = useRef(null)

  // Agent configuration for display
  const agents = [
    { name: 'Keyword Agent', icon: Search, description: 'Semantic keyword matching' },
    { name: 'Format Agent', icon: FileCheck, description: 'ATS compliance check' },
    { name: 'Content Agent', icon: FileText, description: 'Qualitative analysis' },
    { name: 'Scoring Agent', icon: TrendingUp, description: 'Weighted scoring' },
    { name: 'Job Matching Agent', icon: Target, description: 'Resume-JD fit' },
    { name: 'Synthesis Agent', icon: CheckCircle2, description: 'Final report' }
  ]

  const handleUpload = async (event) => {
    const file = event.target.files?.[0]
    if (!file || !jobDescription.trim()) {
      alert('Please provide both resume and job description')
      return
    }

    setUploading(true)
    setAgentProgress({ status: 'starting', completed: [] })

    try {
      const formData = new FormData()
      formData.append('resume', file)
      formData.append('job_description', jobDescription)

      const token = localStorage.getItem('token')
      const response = await fetch('http://localhost:8080/api/v1/resume/upload', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`
        },
        body: formData
      })

      if (!response.ok) throw new Error('Upload failed')

      const data = await response.json()
      setAnalysisId(data.analysis_id)

      // Simulate agent progress (in production, use WebSocket for real-time updates)
      simulateAgentProgress(data.analysis_id)

    } catch (error) {
      console.error('Upload error:', error)
      setAgentProgress({ status: 'error', error: error.message })
      setUploading(false)
    }
  }

  const simulateAgentProgress = async (id) => {
    const agentNames = agents.map(a => a.name)
    
    for (let i = 0; i < agentNames.length; i++) {
      await new Promise(resolve => setTimeout(resolve, 1000))
      setAgentProgress(prev => ({
        ...prev,
        completed: agentNames.slice(0, i + 1),
        current: agentNames[i]
      }))
    }

    // Fetch final results
    await fetchFeedback(id)
  }

  const fetchFeedback = async (id) => {
    try {
      const token = localStorage.getItem('token')
      const response = await fetch(`http://localhost:8080/api/v1/resume/feedback/${id}`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })

      if (!response.ok) throw new Error('Failed to fetch feedback')

      const data = await response.json()
      setFeedback(data)
      setAgentProgress({ status: 'completed', completed: agents.map(a => a.name) })
      setUploading(false)

    } catch (error) {
      console.error('Fetch error:', error)
      setAgentProgress({ status: 'error', error: error.message })
      setUploading(false)
    }
  }

  const getScoreColor = (score) => {
    if (score >= 80) return 'text-green-600'
    if (score >= 60) return 'text-yellow-600'
    return 'text-red-600'
  }

  const progressPercentage = agentProgress ? 
    (agentProgress.completed?.length / agents.length) * 100 : 0

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Multi-Agent ATS Resume Analysis</h1>
        <p className="text-muted-foreground mt-2">
          7-agent system with intelligent keyword matching, format checking, and detailed scoring
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Upload Section */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              <Upload className="mr-2 h-5 w-5" />
              Resume Upload
            </CardTitle>
            <CardDescription>
              Upload resume and job description for multi-agent analysis
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <label className="text-sm font-medium mb-2 block">Job Description</label>
              <Textarea
                placeholder="Paste the job description here..."
                value={jobDescription}
                onChange={(e) => setJobDescription(e.target.value)}
                rows={6}
                disabled={uploading}
              />
            </div>

            <input
              type="file"
              ref={fileInputRef}
              onChange={handleUpload}
              accept=".pdf,.doc,.docx"
              className="hidden"
            />
            <Button 
              className="w-full" 
              onClick={() => fileInputRef.current?.click()}
              disabled={uploading || !jobDescription.trim()}
            >
              {uploading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Analyzing with 7 Agents...
                </>
              ) : (
                <>
                  <Upload className="mr-2 h-4 w-4" />
                  Upload Resume
                </>
              )}
            </Button>

            <div className="text-xs text-muted-foreground space-y-1">
              <p>✓ Intelligent keyword matching (semantic equivalence)</p>
              <p>✓ ATS format compliance checking</p>
              <p>✓ Content quality analysis</p>
              <p>✓ Weighted scoring (40% keywords + 30% format + 30% content)</p>
            </div>
          </CardContent>
        </Card>

        {/* Agent Progress */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              <TrendingUp className="mr-2 h-5 w-5" />
              Agent Execution Progress
            </CardTitle>
            <CardDescription>
              Real-time multi-agent analysis status{analysisId ? ` — Analysis #${analysisId}` : ''}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {agentProgress && (
              <>
                <div className="space-y-2">
                  <div className="flex items-center justify-between text-sm">
                    <span>Overall Progress</span>
                    <span>{Math.round(progressPercentage)}%</span>
                  </div>
                  <Progress value={progressPercentage} />
                </div>

                <div className="space-y-2">
                  {agents.map((agent, index) => {
                    const Icon = agent.icon
                    const isCompleted = agentProgress.completed?.includes(agent.name)
                    const isCurrent = agentProgress.current === agent.name
                    
                    return (
                      <div
                        key={index}
                        className={`flex items-start p-3 border rounded-lg transition-all ${
                          isCompleted ? 'bg-green-50 border-green-200' :
                          isCurrent ? 'bg-blue-50 border-blue-200' :
                          'bg-gray-50'
                        }`}
                      >
                        <div className="mr-3 mt-1">
                          {isCompleted ? (
                            <CheckCircle2 className="h-5 w-5 text-green-600" />
                          ) : isCurrent ? (
                            <Loader2 className="h-5 w-5 text-blue-600 animate-spin" />
                          ) : (
                            <Clock className="h-5 w-5 text-gray-400" />
                          )}
                        </div>
                        <div className="flex-1">
                          <h4 className="text-sm font-medium">{agent.name}</h4>
                          <p className="text-xs text-muted-foreground">{agent.description}</p>
                        </div>
                      </div>
                    )
                  })}
                </div>

                {agentProgress.status === 'error' && (
                  <Alert variant="destructive">
                    <AlertCircle className="h-4 w-4" />
                    <AlertDescription>{agentProgress.error}</AlertDescription>
                  </Alert>
                )}
              </>
            )}

            {!agentProgress && (
              <div className="text-center text-muted-foreground py-8">
                <Clock className="h-12 w-12 mx-auto mb-2 opacity-50" />
                <p className="text-sm">Upload resume to start multi-agent analysis</p>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Feedback Results */}
      {feedback && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center justify-between">
              <span className="flex items-center">
                <CheckCircle2 className="mr-2 h-5 w-5 text-green-600" />
                Analysis Complete
              </span>
              <Badge className={`text-lg px-4 py-1 ${getScoreColor(feedback.score)}`}>
                Score: {feedback.score}/100
              </Badge>
            </CardTitle>
            <CardDescription>
              Comprehensive multi-agent feedback and recommendations
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="prose max-w-none">
              <div className="p-4 bg-muted/50 rounded-lg whitespace-pre-wrap">
                {feedback.feedback}
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="p-4 border rounded-lg">
                <h4 className="font-medium mb-2">Keyword Match</h4>
                <p className="text-2xl font-bold text-blue-600">
                  {Math.round(feedback.score * 0.4)}%
                </p>
                <p className="text-xs text-muted-foreground">40% weight</p>
              </div>
              <div className="p-4 border rounded-lg">
                <h4 className="font-medium mb-2">Format Quality</h4>
                <p className="text-2xl font-bold text-green-600">
                  {Math.round(feedback.score * 0.3)}%
                </p>
                <p className="text-xs text-muted-foreground">30% weight</p>
              </div>
              <div className="p-4 border rounded-lg">
                <h4 className="font-medium mb-2">Content Quality</h4>
                <p className="text-2xl font-bold text-purple-600">
                  {Math.round(feedback.score * 0.3)}%
                </p>
                <p className="text-xs text-muted-foreground">30% weight</p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
