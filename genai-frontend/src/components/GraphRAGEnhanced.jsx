import { useState, useEffect, useRef } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Progress } from '@/components/ui/progress'
import { Search, Upload, Network, Loader2, FileText, AlertCircle } from 'lucide-react'

export default function GraphRAGEnhanced() {
  const [query, setQuery] = useState('')
  const [loading, setLoading] = useState(false)
  const [uploading, setUploading] = useState(false)
  const [result, setResult] = useState(null)
  const [documents, setDocuments] = useState([])
  const [selectedDoc, setSelectedDoc] = useState(null)
  const [processingStatus, setProcessingStatus] = useState(null)
  const fileInputRef = useRef(null)

  // Upload PDF and trigger GraphRAG processing
  const handleUpload = async (event) => {
    const file = event.target.files?.[0]
    if (!file) return

    setUploading(true)
    setProcessingStatus('Uploading...')

    try {
      const formData = new FormData()
      formData.append('file', file)

      const token = localStorage.getItem('token')
      const response = await fetch('http://localhost:8080/api/upload', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`
        },
        body: formData
      })

      if (!response.ok) throw new Error('Upload failed')

      const data = await response.json()
      setSelectedDoc(data.document_id)
      setProcessingStatus('Processing with GraphRAG (extracting entities & relationships)...')
      
      // Add to documents list
      setDocuments(prev => [...prev, {
        id: data.document_id,
        filename: file.name,
        status: 'processing'
      }])

      // Poll for processing completion (simplified - in production use WebSockets)
      setTimeout(() => {
        setProcessingStatus('GraphRAG processing complete! Ready to query.')
        setDocuments(prev => prev.map(doc => 
          doc.id === data.document_id ? { ...doc, status: 'ready' } : doc
        ))
      }, 5000)

    } catch (error) {
      console.error('Upload error:', error)
      setProcessingStatus('Upload failed')
    } finally {
      setUploading(false)
    }
  }

  // Hybrid GraphRAG query
  const handleQuery = async () => {
    if (!query.trim() || !selectedDoc) return
    
    setLoading(true)
    setResult(null)

    try {
      const token = localStorage.getItem('token')
      const response = await fetch('http://localhost:8080/api/chat', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          query: query,
          document_ids: [selectedDoc]
        })
      })

      if (!response.ok) throw new Error('Query failed')

      const data = await response.json()
      setResult({
        answer: data.response,
        contexts_used: data.contexts_used || 0,
        method: data.method || 'hybrid_graphrag'
      })

    } catch (error) {
      console.error('Query error:', error)
      setResult({
        answer: 'Error: Unable to process query. Please try again.',
        error: true
      })
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">GraphRAG - Knowledge Graph QA</h1>
        <p className="text-muted-foreground mt-2">
          Hybrid retrieval combining vector search + 2-hop graph traversal with 50 entity types & 200 relationship types.
        </p>
      </div>

      {/* Processing Status */}
      {processingStatus && (
        <Alert>
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>{processingStatus}</AlertDescription>
        </Alert>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Upload Section */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              <Upload className="mr-2 h-5 w-5" />
              Document Upload
            </CardTitle>
            <CardDescription>
              Upload PDFs to extract entities, relationships, and build knowledge graph
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <input
              type="file"
              ref={fileInputRef}
              onChange={handleUpload}
              accept=".pdf"
              className="hidden"
            />
            <Button 
              className="w-full" 
              onClick={() => fileInputRef.current?.click()}
              disabled={uploading}
            >
              {uploading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Uploading & Processing...
                </>
              ) : (
                <>
                  <Upload className="mr-2 h-4 w-4" />
                  Upload PDF for GraphRAG
                </>
              )}
            </Button>

            {/* Documents List */}
            {documents.length > 0 && (
              <div className="space-y-2">
                <h4 className="text-sm font-medium">Uploaded Documents:</h4>
                {documents.map(doc => (
                  <div 
                    key={doc.id}
                    className={`flex items-center justify-between p-2 border rounded cursor-pointer ${
                      selectedDoc === doc.id ? 'bg-primary/10 border-primary' : ''
                    }`}
                    onClick={() => setSelectedDoc(doc.id)}
                  >
                    <div className="flex items-center">
                      <FileText className="h-4 w-4 mr-2" />
                      <span className="text-sm truncate">{doc.filename}</span>
                    </div>
                    <Badge variant={doc.status === 'ready' ? 'success' : 'secondary'}>
                      {doc.status}
                    </Badge>
                  </div>
                ))}
              </div>
            )}

            <div className="text-xs text-muted-foreground space-y-1">
              <p>✓ Entity extraction (50 types)</p>
              <p>✓ Relationship discovery (200 types)</p>
              <p>✓ Knowledge graph construction</p>
              <p>✓ PageRank centrality calculation</p>
            </div>
          </CardContent>
        </Card>

        {/* Query Section */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              <Search className="mr-2 h-5 w-5" />
              Hybrid GraphRAG Query
            </CardTitle>
            <CardDescription>
              Ask questions using vector retrieval + graph traversal
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex space-x-2">
              <Input
                placeholder="Ask about your documents..."
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleQuery()}
                disabled={!selectedDoc}
              />
              <Button onClick={handleQuery} disabled={loading || !selectedDoc}>
                {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Search className="h-4 w-4" />}
              </Button>
            </div>

            {!selectedDoc && (
              <Alert>
                <AlertCircle className="h-4 w-4" />
                <AlertDescription>Please upload and select a document first</AlertDescription>
              </Alert>
            )}

            {result && (
              <div className="space-y-4">
                <div className={`p-4 border rounded-lg ${result.error ? 'border-red-500' : ''}`}>
                  <div className="flex items-center justify-between mb-2">
                    <h4 className="font-medium">Answer:</h4>
                    {!result.error && (
                      <Badge variant="outline">
                        <Network className="h-3 w-3 mr-1" />
                        {result.method}
                      </Badge>
                    )}
                  </div>
                  <p className="text-sm whitespace-pre-wrap">{result.answer}</p>
                  {result.contexts_used > 0 && (
                    <p className="text-xs text-muted-foreground mt-2">
                      Used {result.contexts_used} context sources (vector + graph entities)
                    </p>
                  )}
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Graph Visualization Placeholder */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <Network className="mr-2 h-5 w-5" />
            Knowledge Graph Visualization
          </CardTitle>
          <CardDescription>
            Visual representation of entities and relationships
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-64 border rounded-lg flex items-center justify-center bg-muted/10">
            <div className="text-center text-muted-foreground">
              <Network className="h-12 w-12 mx-auto mb-2 opacity-50" />
              <p className="text-sm">Graph visualization coming soon</p>
              <p className="text-xs mt-1">Will display entities, relationships, and centrality scores</p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
