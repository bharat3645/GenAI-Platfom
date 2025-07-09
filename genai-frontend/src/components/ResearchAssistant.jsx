import { useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { useAuth } from '@/hooks/useAuth'
import { useToast } from '@/hooks/use-toast'
import { Search, FileText, Loader2, CheckCircle, Clock } from 'lucide-react'

export default function ResearchAssistant() {
  const [query, setQuery] = useState('')
  const [tasks, setTasks] = useState([])
  const [loading, setLoading] = useState(false)
  const { getAuthHeaders } = useAuth()
  const { toast } = useToast()

  const handleSubmitResearch = async () => {
    if (!query.trim()) return

    setLoading(true)

    try {
      const response = await fetch('http://localhost:8080/api/v1/agent/research', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...getAuthHeaders(),
        },
        body: JSON.stringify({ query }),
      })

      if (!response.ok) {
        throw new Error('Failed to submit research task')
      }

      const data = await response.json()
      
      const newTask = {
        id: data.task_id,
        query,
        status: 'started',
        result: null,
        createdAt: new Date().toISOString()
      }

      setTasks(prev => [newTask, ...prev])
      setQuery('')

      toast({
        title: "Research task started",
        description: "Your research task has been submitted and is being processed.",
      })

      // Poll for results
      pollTaskResult(data.task_id)
    } catch (error) {
      toast({
        title: "Failed to submit research",
        description: error.message,
        variant: "destructive",
      })
    } finally {
      setLoading(false)
    }
  }

  const pollTaskResult = async (taskId) => {
    const maxAttempts = 20
    let attempts = 0

    const poll = async () => {
      try {
        const response = await fetch(`http://localhost:8080/api/v1/agent/research/${taskId}`, {
          headers: getAuthHeaders(),
        })

        if (response.ok) {
          const data = await response.json()
          
          setTasks(prev => prev.map(task => 
            task.id === taskId 
              ? { ...task, status: data.status, result: data.result, completedAt: data.completed_at }
              : task
          ))

          if (data.status === 'completed') {
            toast({
              title: "Research completed",
              description: "Your research task has been completed successfully.",
            })
            return
          }
        }

        attempts++
        if (attempts < maxAttempts) {
          setTimeout(poll, 3000) // Poll every 3 seconds
        }
      } catch (error) {
        console.error('Polling error:', error)
      }
    }

    poll()
  }

  const getStatusIcon = (status) => {
    switch (status) {
      case 'completed':
        return <CheckCircle className="h-4 w-4 text-green-600" />
      case 'started':
      case 'pending':
        return <Clock className="h-4 w-4 text-yellow-600" />
      default:
        return <Loader2 className="h-4 w-4 animate-spin text-blue-600" />
    }
  }

  const getStatusColor = (status) => {
    switch (status) {
      case 'completed':
        return 'default'
      case 'started':
      case 'pending':
        return 'secondary'
      default:
        return 'outline'
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Agentic Research Assistant</h1>
        <p className="text-muted-foreground mt-2">
          AI-powered research agent that autonomously gathers and synthesizes information from multiple sources.
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Submit Research */}
        <div className="lg:col-span-1">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <Search className="mr-2 h-5 w-5" />
                New Research
              </CardTitle>
              <CardDescription>
                Submit a research query for AI analysis
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <Textarea
                placeholder="Enter your research question or topic..."
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                className="min-h-[120px]"
              />
              <Button
                onClick={handleSubmitResearch}
                disabled={loading || !query.trim()}
                className="w-full"
              >
                {loading ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Submitting...
                  </>
                ) : (
                  <>
                    <Search className="mr-2 h-4 w-4" />
                    Start Research
                  </>
                )}
              </Button>
            </CardContent>
          </Card>
        </div>

        {/* Research Tasks */}
        <div className="lg:col-span-2">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <FileText className="mr-2 h-5 w-5" />
                Research Tasks
              </CardTitle>
              <CardDescription>
                Track your research tasks and view results
              </CardDescription>
            </CardHeader>
            <CardContent>
              {tasks.length === 0 ? (
                <div className="text-center text-muted-foreground py-8">
                  No research tasks yet. Submit a research query to get started.
                </div>
              ) : (
                <div className="space-y-4">
                  {tasks.map((task) => (
                    <div key={task.id} className="border rounded-lg p-4">
                      <div className="flex items-start justify-between mb-2">
                        <div className="flex items-center space-x-2">
                          {getStatusIcon(task.status)}
                          <Badge variant={getStatusColor(task.status)}>
                            {task.status}
                          </Badge>
                        </div>
                        <span className="text-xs text-muted-foreground">
                          {new Date(task.createdAt).toLocaleString()}
                        </span>
                      </div>
                      
                      <h4 className="font-medium mb-2">Research Query:</h4>
                      <p className="text-sm text-muted-foreground mb-3">{task.query}</p>
                      
                      {task.result && (
                        <div>
                          <h4 className="font-medium mb-2">Results:</h4>
                          <div className="bg-muted/50 p-3 rounded border">
                            <p className="text-sm whitespace-pre-wrap">{task.result}</p>
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

