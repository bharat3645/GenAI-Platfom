# 🚀 Quick Start Guide - GenAI Platform

## ⚡ Fastest Way to Get Started

### Option 1: Using Neon Database (Recommended for Production)

1. **Get your Neon connection string** from [Neon Dashboard](https://console.neon.tech)
   - It looks like: `postgres://username:password@ep-cool-name-123456.us-east-2.aws.neon.tech/dbname`

2. **Update the `.env` file** in `genai-platform/`:
   ```env
   DATABASE_URL=your-neon-connection-string-here
   ```

3. **Run the setup script**:
   ```powershell
   # Windows
   .\setup.ps1
   ```

4. **Run migrations**:
   ```bash
   psql $DATABASE_URL -f genai-platform/migrations/001_initial_schema.sql
   ```

5. **Start the servers**:
   ```powershell
   # Terminal 1 - Backend
   cd genai-platform
   go run cmd/server/main.go

   # Terminal 2 - Frontend
   cd genai-frontend
   npm run dev
   ```

### Option 2: Using Docker (Easiest)

```bash
# Start everything with one command
docker-compose up --build

# Access:
# Frontend: http://localhost:5173
# Backend: http://localhost:8080
```

### Option 3: Local PostgreSQL

1. **Install PostgreSQL** locally

2. **Create database**:
   ```sql
   CREATE DATABASE genai_platform;
   ```

3. **Update `.env`**:
   ```env
   DATABASE_URL=postgres://postgres:yourpassword@localhost:5432/genai_platform
   ```

4. Follow steps 3-5 from Option 1

## 🔧 Configuration Checklist

### Required Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://user:pass@host:5432/db` |
| `JWT_SECRET` | Secret key for JWT tokens | Generate random 32+ chars |
| `OPENAI_API_KEY` | OpenAI API key (optional) | `sk-...` |
| `PORT` | Backend server port | `8080` |

### Optional Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `GEMINI_API_KEY` | Google Gemini API key | - |
| `SENDGRID_API_KEY` | Email service key | - |
| `UPLOAD_DIR` | File upload directory | `./uploads` |

## 📝 Important Notes

### Database Connection String Format

**For Neon/Supabase:**
```
postgres://user:password@host.neon.tech:5432/dbname?sslmode=require
```

**Special Characters in Password:**
- `#` → `%23`
- `@` → `%40`
- `:` → `%3A`
- `/` → `%2F`

**Example:**
- Password: `MyPass#123`
- Encoded: `MyPass%23123`

### Common Issues & Solutions

#### Issue: "no such host" error
**Solution:** Check your DATABASE_URL hostname is correct. Copy it exactly from Neon/Supabase dashboard.

#### Issue: "password authentication failed"
**Solution:** Ensure password special characters are URL-encoded in DATABASE_URL.

#### Issue: Frontend can't reach backend
**Solution:** Make sure backend is running on port 8080 and CORS is configured (already done).

#### Issue: npm install fails
**Solution:** Use `npm install --legacy-peer-deps` flag.

## 🎯 What's Working Now

✅ Frontend development server configured and ready  
✅ Backend API structure complete with all endpoints  
✅ CORS configured for cross-origin requests  
✅ Vite proxy configured for API calls  
✅ Database migrations ready  
✅ Docker configuration complete  
✅ Authentication system implemented  
✅ File upload handling configured  

## 🔜 What You Need to Do

1. **Get a Neon Database**:
   - Sign up at [neon.tech](https://neon.tech)
   - Create a new project
   - Copy the connection string

2. **Update Configuration**:
   - Paste connection string in `genai-platform/.env`
   - Add your OpenAI API key if you have one

3. **Run Migrations**:
   - Execute the SQL migration file

4. **Start the Application**:
   - Run backend and frontend as shown above

## 📚 Next Steps After Setup

1. Open browser to `http://localhost:5173`
2. Register a new account
3. Try uploading a PDF in the PDF Chat module
4. Test other features (Resume Analyzer, Research Assistant, etc.)

## 🆘 Getting Help

If you encounter issues:

1. Check the logs in the terminal
2. Verify your `.env` file configuration
3. Ensure database is accessible
4. Review the main [README.md](README.md) for detailed documentation

## 📊 System Requirements

- **Go**: 1.21 or higher
- **Node.js**: 18 or higher
- **PostgreSQL**: 14 or higher (or Neon account)
- **RAM**: 4GB minimum
- **Disk**: 2GB free space

## 🎉 Success Indicators

You'll know everything is working when you see:

**Backend Terminal:**
```
2025/11/05 14:40:00 No .env file found, using environment variables
2025/11/05 14:40:01 Server starting on port 8080
```

**Frontend Terminal:**
```
VITE v6.3.5  ready in 661 ms

➜  Local:   http://localhost:5173/
➜  Network: use --host to expose
```

**Browser:**
- Login page loads at `http://localhost:5173`
- No console errors
- Can register and login successfully

---

**You're almost there! Just need to connect the database and you're ready to go! 🚀**
