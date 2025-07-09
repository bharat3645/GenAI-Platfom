import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { useNavigate } from 'react-router-dom'
import { 
  MessageSquare, 
  Search, 
  FileText, 
  FileUser, 
  Database,
  Brain,
  Zap,
  Shield,
  Globe
} from 'lucide-react'

const features = [
  {
    name: 'Multi-PDF Chat',
    description: 'Upload multiple PDFs and chat with your documents using advanced RAG technology.',
    icon: MessageSquare,
    href: '/dashboard/pdf-chat',
    color: 'text-blue-600'
  },
  {
    name: 'GraphRAG QA',
    description: 'Extract entities and relationships to build knowledge graphs for enhanced Q&A.',
    icon: Search,
    href: '/dashboard/graph-rag',
    color: 'text-green-600'
  },
  {
    name: 'Research Assistant',
    description: 'AI-powered research agent that autonomously gathers and synthesizes information.',
    icon: FileText,
    href: '/dashboard/research',
    color: 'text-purple-600'
  },
  {
    name: 'Resume Feedback',
    description: 'Get AI-powered feedback and ATS scoring for your resume optimization.',
    icon: FileUser,
    href: '/dashboard/resume',
    color: 'text-orange-600'
  },
  {
    name: 'Text to SQL',
    description: 'Convert natural language queries into SQL and execute them on your database.',
    icon: Database,
    href: '/dashboard/sql',
    color: 'text-red-600'
  }
]

const stats = [
  { name: 'AI Models Integrated', value: '3+', icon: Brain },
  { name: 'Processing Speed', value: '10x', icon: Zap },
  { name: 'Data Security', value: '100%', icon: Shield },
  { name: 'Global Access', value: '24/7', icon: Globe },
]

export default function DashboardHome() {
  const navigate = useNavigate()

  return (
    <div className="space-y-8">
      {/* Welcome Section */}
      <div className="text-center">
        <h1 className="text-4xl font-bold tracking-tight text-foreground sm:text-6xl">
          Welcome to GenAI Platform
        </h1>
        <p className="mt-6 text-lg leading-8 text-muted-foreground max-w-3xl mx-auto">
          An advanced, modular AI-powered platform enabling sophisticated Retrieval-Augmented Generation (RAG), 
          agentic AI workflows, and automated NLP tasks. Built for researchers, professionals, and knowledge workers.
        </p>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => (
          <Card key={stat.name}>
            <CardContent className="p-6">
              <div className="flex items-center">
                <stat.icon className="h-8 w-8 text-primary" />
                <div className="ml-4">
                  <p className="text-2xl font-bold text-foreground">{stat.value}</p>
                  <p className="text-sm text-muted-foreground">{stat.name}</p>
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Features Grid */}
      <div>
        <h2 className="text-2xl font-bold text-foreground mb-6">Platform Features</h2>
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {features.map((feature) => (
            <Card key={feature.name} className="hover:shadow-lg transition-shadow cursor-pointer">
              <CardHeader>
                <div className="flex items-center space-x-3">
                  <feature.icon className={`h-8 w-8 ${feature.color}`} />
                  <CardTitle className="text-lg">{feature.name}</CardTitle>
                </div>
              </CardHeader>
              <CardContent>
                <CardDescription className="text-sm mb-4">
                  {feature.description}
                </CardDescription>
                <Button 
                  onClick={() => navigate(feature.href)}
                  className="w-full"
                  variant="outline"
                >
                  Get Started
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>

      {/* Technology Stack */}
      <Card>
        <CardHeader>
          <CardTitle>Technology Stack</CardTitle>
          <CardDescription>
            Built with cutting-edge technologies for optimal performance and scalability
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">
            <div className="p-4 border rounded-lg">
              <h3 className="font-semibold">Frontend</h3>
              <p className="text-sm text-muted-foreground">React, Tailwind CSS</p>
            </div>
            <div className="p-4 border rounded-lg">
              <h3 className="font-semibold">Backend</h3>
              <p className="text-sm text-muted-foreground">Go, Chi Router</p>
            </div>
            <div className="p-4 border rounded-lg">
              <h3 className="font-semibold">AI Models</h3>
              <p className="text-sm text-muted-foreground">GPT-4, Gemini Pro</p>
            </div>
            <div className="p-4 border rounded-lg">
              <h3 className="font-semibold">Databases</h3>
              <p className="text-sm text-muted-foreground">PostgreSQL, FAISS, Neo4j</p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

