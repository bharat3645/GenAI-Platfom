import { useState, useRef } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { useAuth } from '@/hooks/useAuth'
import { useToast } from '@/hooks/use-toast'
import { Upload, MessageSquare, FileText, Loader2, Send } from 'lucide-react'

export default function PDFChat() {
  const [files, setFiles] = useState([])
  const [uploading, setUploading] = useState(false)
  const [messages, setMessages] = useState([])
  const [query, setQuery] = useState('')
  const [loading, setLoading] = useState(false)
  const [sessionId, setSessionId] = useState(null)
  const fileInputRef = useRef(null)
  const { getAuthHeaders } = useAuth()
  const { toast } = useToast()

  const handleFileUpload = async (event) => {
    const selectedFiles = Array.from(event.target.files)
    if (selectedFiles.length === 0) return

    setUploading(true)

    try {
      const uploadPromises = selectedFiles.map(async (file) => {
        const formData = new FormData()
        formData.append('file', file)

        const response = await fetch('http://localhost:8080/api/v1/pdf/upload', {
          method: 'POST',
          headers: getAuthHeaders(),
          body: formData,
        })

        if (!response.ok) {
          throw new Error(`Failed to upload ${file.name}`)
        }

        const data = await response.json()
        return {
          id: data.document_id,
          name: file.name,
          status: data.status
        }
      })

      const uploadedFiles = await Promise.all(uploadPromises)
      setFiles(prev => [...prev, ...uploadedFiles])

      toast({
        title: "Files uploaded successfully",
        description: `${uploadedFiles.length} file(s) have been uploaded and are being processed.`,
      })
    } catch (error) {
      toast({
        title: "Upload failed",
        description: error.message,
        variant: "destructive",
      })
    } finally {
      setUploading(false)
    }
  }

  const handleSendMessage = async () => {
    if (!query.trim() || files.length === 0) return

    const userMessage = { role: 'user', content: query }
    setMessages(prev => [...prev, userMessage])
    setLoading(true)

    try {
      const response = await fetch('http://localhost:8080/api/v1/chat/query', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...getAuthHeaders(),
        },
        body: JSON.stringify({
          query,
          document_ids: files.map(f => f.id),
          session_id: sessionId
        }),
      })

      if (!response.ok) {
        throw new Error('Failed to get response')
      }

      const data = await response.json()
      
      if (!sessionId) {
        setSessionId(data.session_id)
      }

      const assistantMessage = { 
        role: 'assistant', 
        content: data.response,
        context: data.context
      }
      setMessages(prev => [...prev, assistantMessage])
      setQuery('')
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

  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSendMessage()
    }
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Multi-PDF Chat</h1>
        <p className="text-muted-foreground mt-2">
          Upload multiple PDFs and chat with your documents using advanced RAG technology.
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* File Upload Section */}
        <div className="lg:col-span-1">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <FileText className="mr-2 h-5 w-5" />
                Documents
              </CardTitle>
              <CardDescription>
                Upload PDF files to chat with their content
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <Button
                onClick={() => fileInputRef.current?.click()}
                disabled={uploading}
                className="w-full"
                variant="outline"
              >
                {uploading ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Uploading...
                  </>
                ) : (
                  <>
                    <Upload className="mr-2 h-4 w-4" />
                    Upload PDFs
                  </>
                )}
              </Button>
              
              <input
                ref={fileInputRef}
                type="file"
                multiple
                accept=".pdf"
                onChange={handleFileUpload}
                className="hidden"
              />

              {files.length > 0 && (
                <div className="space-y-2">
                  <h4 className="text-sm font-medium">Uploaded Files:</h4>
                  {files.map((file, index) => (
                    <div key={index} className="flex items-center justify-between p-2 border rounded">
                      <span className="text-sm truncate">{file.name}</span>
                      <Badge variant={file.status === 'uploaded' ? 'default' : 'secondary'}>
                        {file.status}
                      </Badge>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </div>

        {/* Chat Section */}
        <div className="lg:col-span-2">
          <Card className="h-[600px] flex flex-col">
            <CardHeader>
              <CardTitle className="flex items-center">
                <MessageSquare className="mr-2 h-5 w-5" />
                Chat with Documents
              </CardTitle>
            </CardHeader>
            <CardContent className="flex-1 flex flex-col">
              {/* Messages */}
              <div className="flex-1 overflow-y-auto space-y-4 mb-4 p-4 border rounded-lg bg-muted/20">
                {messages.length === 0 ? (
                  <div className="text-center text-muted-foreground py-8">
                    {files.length === 0 
                      ? "Upload some PDF files to start chatting"
                      : "Ask a question about your uploaded documents"
                    }
                  </div>
                ) : (
                  messages.map((message, index) => (
                    <div
                      key={index}
                      className={`flex ${message.role === 'user' ? 'justify-end' : 'justify-start'}`}
                    >
                      <div
                        className={`max-w-[80%] p-3 rounded-lg ${
                          message.role === 'user'
                            ? 'bg-primary text-primary-foreground'
                            : 'bg-card border'
                        }`}
                      >
                        <p className="whitespace-pre-wrap">{message.content}</p>
                        {message.context && (
                          <details className="mt-2 text-xs opacity-70">
                            <summary className="cursor-pointer">Context used</summary>
                            <p className="mt-1 whitespace-pre-wrap">{message.context}</p>
                          </details>
                        )}
                      </div>
                    </div>
                  ))
                )}
                {loading && (
                  <div className="flex justify-start">
                    <div className="bg-card border p-3 rounded-lg">
                      <Loader2 className="h-4 w-4 animate-spin" />
                    </div>
                  </div>
                )}
              </div>

              {/* Input */}
              <div className="flex space-x-2">
                <Textarea
                  placeholder="Ask a question about your documents..."
                  value={query}
                  onChange={(e) => setQuery(e.target.value)}
                  onKeyPress={handleKeyPress}
                  disabled={loading || files.length === 0}
                  className="flex-1 min-h-[60px] max-h-[120px]"
                />
                <Button
                  onClick={handleSendMessage}
                  disabled={loading || !query.trim() || files.length === 0}
                  size="lg"
                >
                  <Send className="h-4 w-4" />
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}

