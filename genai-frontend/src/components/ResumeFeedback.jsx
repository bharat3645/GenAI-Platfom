import { useState, useRef } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { useAuth } from '@/hooks/useAuth'
import { useToast } from '@/hooks/use-toast'
import { Upload, FileUser, Loader2, CheckCircle, Star } from 'lucide-react'

export default function ResumeFeedback() {
  const [analyses, setAnalyses] = useState([])
  const [uploading, setUploading] = useState(false)
  const [jobDescription, setJobDescription] = useState('')
  const fileInputRef = useRef(null)
  const { getAuthHeaders } = useAuth()
  const { toast } = useToast()

  const handleResumeUpload = async (event) => {
    const file = event.target.files[0]
    if (!file) return

    setUploading(true)

    try {
      const formData = new FormData()
      formData.append('resume', file)
      formData.append('job_description', jobDescription)

      const response = await fetch('http://localhost:8080/api/v1/resume/upload', {
        method: 'POST',
        headers: getAuthHeaders(),
        body: formData,
      })

      if (!response.ok) {
        throw new Error('Failed to upload resume')
      }

      const data = await response.json()
      
      const newAnalysis = {
        id: data.analysis_id,
        filename: file.name,
        jobDescription,
        status: 'processing',
        feedback: null,
        score: null,
        createdAt: new Date().toISOString()
      }

      setAnalyses(prev => [newAnalysis, ...prev])
      setJobDescription('')

      toast({
        title: "Resume uploaded successfully",
        description: "Your resume is being analyzed. Results will appear shortly.",
      })

      // Poll for results
      pollAnalysisResult(data.analysis_id)
    } catch (error) {
      toast({
        title: "Upload failed",
        description: error.message,
        variant: "destructive",
      })
    } finally {
      setUploading(false)
      if (fileInputRef.current) {
        fileInputRef.current.value = ''
      }
    }
  }

  const pollAnalysisResult = async (analysisId) => {
    const maxAttempts = 15
    let attempts = 0

    const poll = async () => {
      try {
        const response = await fetch(`http://localhost:8080/api/v1/resume/feedback/${analysisId}`, {
          headers: getAuthHeaders(),
        })

        if (response.ok) {
          const data = await response.json()
          
          setAnalyses(prev => prev.map(analysis => 
            analysis.id === analysisId 
              ? { 
                  ...analysis, 
                  status: data.status, 
                  feedback: data.feedback, 
                  score: data.score,
                  completedAt: data.completed_at 
                }
              : analysis
          ))

          if (data.status === 'completed') {
            toast({
              title: "Analysis completed",
              description: `Your resume has been analyzed with a score of ${data.score}/100.`,
            })
            return
          }
        }

        attempts++
        if (attempts < maxAttempts) {
          setTimeout(poll, 2000) // Poll every 2 seconds
        }
      } catch (error) {
        console.error('Polling error:', error)
      }
    }

    poll()
  }

  const getScoreColor = (score) => {
    if (score >= 80) return 'text-green-600'
    if (score >= 60) return 'text-yellow-600'
    return 'text-red-600'
  }

  const getScoreLabel = (score) => {
    if (score >= 80) return 'Excellent'
    if (score >= 60) return 'Good'
    return 'Needs Improvement'
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Resume Feedback Bot</h1>
        <p className="text-muted-foreground mt-2">
          Get AI-powered feedback and ATS scoring for your resume optimization.
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Upload Section */}
        <div className="lg:col-span-1">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <Upload className="mr-2 h-5 w-5" />
                Upload Resume
              </CardTitle>
              <CardDescription>
                Upload your resume and job description for analysis
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <Textarea
                placeholder="Paste the job description here (optional but recommended for better analysis)..."
                value={jobDescription}
                onChange={(e) => setJobDescription(e.target.value)}
                className="min-h-[120px]"
              />
              
              <Button
                onClick={() => fileInputRef.current?.click()}
                disabled={uploading}
                className="w-full"
                variant="outline"
              >
                {uploading ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Analyzing...
                  </>
                ) : (
                  <>
                    <Upload className="mr-2 h-4 w-4" />
                    Upload Resume
                  </>
                )}
              </Button>
              
              <input
                ref={fileInputRef}
                type="file"
                accept=".pdf,.doc,.docx"
                onChange={handleResumeUpload}
                className="hidden"
              />

              <p className="text-xs text-muted-foreground text-center">
                Supports PDF, DOC, and DOCX files
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Analysis Results */}
        <div className="lg:col-span-2">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <FileUser className="mr-2 h-5 w-5" />
                Analysis Results
              </CardTitle>
              <CardDescription>
                View your resume analysis and feedback
              </CardDescription>
            </CardHeader>
            <CardContent>
              {analyses.length === 0 ? (
                <div className="text-center text-muted-foreground py-8">
                  No resume analyses yet. Upload a resume to get started.
                </div>
              ) : (
                <div className="space-y-6">
                  {analyses.map((analysis) => (
                    <div key={analysis.id} className="border rounded-lg p-4">
                      <div className="flex items-start justify-between mb-4">
                        <div>
                          <h4 className="font-medium">{analysis.filename}</h4>
                          <p className="text-xs text-muted-foreground">
                            {new Date(analysis.createdAt).toLocaleString()}
                          </p>
                        </div>
                        <div className="flex items-center space-x-2">
                          {analysis.status === 'completed' ? (
                            <CheckCircle className="h-4 w-4 text-green-600" />
                          ) : (
                            <Loader2 className="h-4 w-4 animate-spin text-blue-600" />
                          )}
                          <Badge variant={analysis.status === 'completed' ? 'default' : 'secondary'}>
                            {analysis.status}
                          </Badge>
                        </div>
                      </div>

                      {analysis.score !== null && (
                        <div className="mb-4">
                          <div className="flex items-center justify-between mb-2">
                            <span className="text-sm font-medium">ATS Score</span>
                            <div className="flex items-center space-x-2">
                              <span className={`text-lg font-bold ${getScoreColor(analysis.score)}`}>
                                {analysis.score}/100
                              </span>
                              <Badge variant="outline" className={getScoreColor(analysis.score)}>
                                {getScoreLabel(analysis.score)}
                              </Badge>
                            </div>
                          </div>
                          <Progress value={analysis.score} className="h-2" />
                        </div>
                      )}

                      {analysis.jobDescription && (
                        <div className="mb-4">
                          <h5 className="font-medium mb-2">Job Description Used:</h5>
                          <div className="bg-muted/50 p-3 rounded border text-sm">
                            {analysis.jobDescription.length > 200 
                              ? `${analysis.jobDescription.substring(0, 200)}...`
                              : analysis.jobDescription
                            }
                          </div>
                        </div>
                      )}

                      {analysis.feedback && (
                        <div>
                          <h5 className="font-medium mb-2 flex items-center">
                            <Star className="mr-2 h-4 w-4" />
                            Feedback & Recommendations:
                          </h5>
                          <div className="bg-muted/50 p-3 rounded border">
                            <p className="text-sm whitespace-pre-wrap">{analysis.feedback}</p>
                          </div>
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}

