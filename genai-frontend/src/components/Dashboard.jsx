import { useState } from 'react'
import { Routes, Route, useNavigate, useLocation } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { useAuth } from '@/hooks/useAuth'
import { 
  Brain, 
  FileText, 
  MessageSquare, 
  Search, 
  FileUser, 
  Database,
  LogOut,
  Menu,
  X
} from 'lucide-react'
import PDFChat from '@/components/PDFChat'
import GraphRAG from '@/components/GraphRAG'
import ResearchAssistant from '@/components/ResearchAssistant'
import ResumeFeedback from '@/components/ResumeFeedback'
import TextToSQL from '@/components/TextToSQL'
import DashboardHome from '@/components/DashboardHome'

const navigation = [
  { name: 'Overview', href: '/dashboard', icon: Brain },
  { name: 'PDF Chat', href: '/dashboard/pdf-chat', icon: MessageSquare },
  { name: 'GraphRAG', href: '/dashboard/graph-rag', icon: Search },
  { name: 'Research Assistant', href: '/dashboard/research', icon: FileText },
  { name: 'Resume Feedback', href: '/dashboard/resume', icon: FileUser },
  { name: 'Text to SQL', href: '/dashboard/sql', icon: Database },
]

export default function Dashboard() {
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const { user, logout } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  const isActive = (href) => {
    if (href === '/dashboard') {
      return location.pathname === '/dashboard'
    }
    return location.pathname.startsWith(href)
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Mobile sidebar */}
      <div className={`fixed inset-0 z-50 lg:hidden ${sidebarOpen ? 'block' : 'hidden'}`}>
        <div className="fixed inset-0 bg-black/50" onClick={() => setSidebarOpen(false)} />
        <div className="fixed left-0 top-0 h-full w-64 bg-card border-r">
          <SidebarContent 
            navigation={navigation} 
            user={user} 
            onLogout={handleLogout}
            isActive={isActive}
            onClose={() => setSidebarOpen(false)}
          />
        </div>
      </div>

      {/* Desktop sidebar */}
      <div className="hidden lg:fixed lg:inset-y-0 lg:flex lg:w-64 lg:flex-col">
        <div className="flex flex-col flex-grow bg-card border-r">
          <SidebarContent 
            navigation={navigation} 
            user={user} 
            onLogout={handleLogout}
            isActive={isActive}
          />
        </div>
      </div>

      {/* Main content */}
      <div className="lg:pl-64">
        {/* Top bar */}
        <div className="sticky top-0 z-40 flex h-16 shrink-0 items-center gap-x-4 border-b bg-card px-4 shadow-sm sm:gap-x-6 sm:px-6 lg:px-8">
          <Button
            variant="ghost"
            size="sm"
            className="lg:hidden"
            onClick={() => setSidebarOpen(true)}
          >
            <Menu className="h-6 w-6" />
          </Button>
          
          <div className="flex flex-1 gap-x-4 self-stretch lg:gap-x-6">
            <div className="flex flex-1 items-center">
              <h1 className="text-lg font-semibold">
                {navigation.find(item => isActive(item.href))?.name || 'GenAI Platform'}
              </h1>
            </div>
            <div className="flex items-center gap-x-4 lg:gap-x-6">
              <span className="text-sm text-muted-foreground">
                Welcome, {user?.email}
              </span>
            </div>
          </div>
        </div>

        {/* Page content */}
        <main className="py-6">
          <div className="px-4 sm:px-6 lg:px-8">
            <Routes>
              <Route path="/" element={<DashboardHome />} />
              <Route path="/pdf-chat" element={<PDFChat />} />
              <Route path="/graph-rag" element={<GraphRAG />} />
              <Route path="/research" element={<ResearchAssistant />} />
              <Route path="/resume" element={<ResumeFeedback />} />
              <Route path="/sql" element={<TextToSQL />} />
            </Routes>
          </div>
        </main>
      </div>
    </div>
  )
}

function SidebarContent({ navigation, user, onLogout, isActive, onClose }) {
  const navigate = useNavigate()

  const handleNavigation = (href) => {
    navigate(href)
    if (onClose) onClose()
  }

  return (
    <>
      <div className="flex h-16 shrink-0 items-center px-6 border-b">
        <div className="flex items-center">
          <Brain className="h-8 w-8 text-primary" />
          <span className="ml-2 text-xl font-bold">GenAI Platform</span>
        </div>
        {onClose && (
          <Button
            variant="ghost"
            size="sm"
            className="ml-auto lg:hidden"
            onClick={onClose}
          >
            <X className="h-6 w-6" />
          </Button>
        )}
      </div>
      
      <nav className="flex flex-1 flex-col px-6 py-4">
        <ul role="list" className="flex flex-1 flex-col gap-y-7">
          <li>
            <ul role="list" className="-mx-2 space-y-1">
              {navigation.map((item) => (
                <li key={item.name}>
                  <Button
                    variant={isActive(item.href) ? "secondary" : "ghost"}
                    className="w-full justify-start"
                    onClick={() => handleNavigation(item.href)}
                  >
                    <item.icon className="mr-3 h-5 w-5" />
                    {item.name}
                  </Button>
                </li>
              ))}
            </ul>
          </li>
          
          <li className="mt-auto">
            <div className="border-t pt-4">
              <div className="px-2 py-2 text-sm text-muted-foreground">
                {user?.email}
              </div>
              <Button
                variant="ghost"
                className="w-full justify-start text-destructive hover:text-destructive"
                onClick={onLogout}
              >
                <LogOut className="mr-3 h-5 w-5" />
                Sign out
              </Button>
            </div>
          </li>
        </ul>
      </nav>
    </>
  )
}

