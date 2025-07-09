import { useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { useAuth } from '@/hooks/useAuth'
import { useToast } from '@/hooks/use-toast'
import { Database, Play, Code, Loader2, Table } from 'lucide-react'

export default function TextToSQL() {
  const [query, setQuery] = useState('')
  const [queries, setQueries] = useState([])
  const [loading, setLoading] = useState(false)
  const { getAuthHeaders } = useAuth()
  const { toast } = useToast()

  const handleExecuteQuery = async () => {
    if (!query.trim()) return

    setLoading(true)

    try {
      const response = await fetch('http://localhost:8080/api/v1/sql/query', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...getAuthHeaders(),
        },
        body: JSON.stringify({ query }),
      })

      if (!response.ok) {
        throw new Error('Failed to execute query')
      }

      const data = await response.json()
      
      const newQuery = {
        id: data.query_id,
        naturalQuery: query,
        generatedSQL: data.sql,
        resultData: data.result_data,
        createdAt: new Date().toISOString()
      }

      setQueries(prev => [newQuery, ...prev])
      setQuery('')

      toast({
        title: "Query executed successfully",
        description: "Your natural language query has been converted to SQL and executed.",
      })
    } catch (error) {
      toast({
        title: "Query failed",
        description: error.message,
        variant: "destructive",
      })
    } finally {
      setLoading(false)
    }
  }

  const exampleQueries = [
    "Show me all users who registered in the last 30 days",
    "Find the top 5 most uploaded document types",
    "Count how many chat sessions were created this month",
    "List all users with their total number of documents"
  ]

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Text to SQL</h1>
        <p className="text-muted-foreground mt-2">
          Convert natural language queries into SQL and execute them on your database.
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Query Input */}
        <div className="lg:col-span-1">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <Database className="mr-2 h-5 w-5" />
                Natural Language Query
              </CardTitle>
              <CardDescription>
                Describe what data you want to retrieve
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <Textarea
                placeholder="Describe what data you want to retrieve in plain English..."
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                className="min-h-[120px]"
              />
              
              <Button
                onClick={handleExecuteQuery}
                disabled={loading || !query.trim()}
                className="w-full"
              >
                {loading ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Executing...
                  </>
                ) : (
                  <>
                    <Play className="mr-2 h-4 w-4" />
                    Execute Query
                  </>
                )}
              </Button>

              <div className="space-y-2">
                <h4 className="text-sm font-medium">Example queries:</h4>
                {exampleQueries.map((example, index) => (
                  <Button
                    key={index}
                    variant="outline"
                    size="sm"
                    className="w-full text-left justify-start h-auto p-2 text-xs"
                    onClick={() => setQuery(example)}
                  >
                    {example}
                  </Button>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Query Results */}
        <div className="lg:col-span-2">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <Table className="mr-2 h-5 w-5" />
                Query Results
              </CardTitle>
              <CardDescription>
                Generated SQL and execution results
              </CardDescription>
            </CardHeader>
            <CardContent>
              {queries.length === 0 ? (
                <div className="text-center text-muted-foreground py-8">
                  No queries executed yet. Enter a natural language query to get started.
                </div>
              ) : (
                <div className="space-y-6">
                  {queries.map((queryResult) => (
                    <div key={queryResult.id} className="border rounded-lg p-4">
                      <div className="flex items-start justify-between mb-4">
                        <Badge variant="outline">Query #{queryResult.id}</Badge>
                        <span className="text-xs text-muted-foreground">
                          {new Date(queryResult.createdAt).toLocaleString()}
                        </span>
                      </div>

                      <div className="space-y-4">
                        <div>
                          <h4 className="font-medium mb-2">Natural Language Query:</h4>
                          <p className="text-sm bg-muted/50 p-3 rounded border">
                            {queryResult.naturalQuery}
                          </p>
                        </div>

                        <div>
                          <h4 className="font-medium mb-2 flex items-center">
                            <Code className="mr-2 h-4 w-4" />
                            Generated SQL:
                          </h4>
                          <pre className="text-xs bg-black text-green-400 p-3 rounded border overflow-x-auto">
                            {queryResult.generatedSQL}
                          </pre>
                        </div>

                        <div>
                          <h4 className="font-medium mb-2">Results:</h4>
                          <div className="bg-muted/50 p-3 rounded border">
                            {queryResult.resultData.message ? (
                              <p className="text-sm">{queryResult.resultData.message}</p>
                            ) : (
                              <div className="text-sm">
                                <p className="text-muted-foreground">
                                  In a production environment, this would show the actual query results in a table format.
                                </p>
                              </div>
                            )}
                          </div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Database Schema Info */}
      <Card>
        <CardHeader>
          <CardTitle>Available Tables</CardTitle>
          <CardDescription>
            Current database schema for reference
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {[
              { name: 'users', description: 'User accounts and authentication' },
              { name: 'documents', description: 'Uploaded PDF documents' },
              { name: 'chat_sessions', description: 'PDF chat conversations' },
              { name: 'chat_messages', description: 'Individual chat messages' },
              { name: 'research_tasks', description: 'Research assistant tasks' },
              { name: 'resume_analyses', description: 'Resume feedback results' },
              { name: 'sql_queries', description: 'Text-to-SQL query history' }
            ].map((table) => (
              <div key={table.name} className="p-3 border rounded-lg">
                <h4 className="font-medium text-sm">{table.name}</h4>
                <p className="text-xs text-muted-foreground mt-1">{table.description}</p>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

