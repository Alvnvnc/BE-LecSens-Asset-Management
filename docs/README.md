# Documentation Index

Welcome to the LecSens Asset Management Backend documentation. This directory contains comprehensive guides and references for developers and administrators.

## 📋 Table of Contents

### Getting Started
- **[README](../README.md)** - Main project overview and setup instructions
- **[Quick Reference](CLI_Quick_Reference.md)** - Essential command line commands

### Command Line Tools
- **[Command Line Tools Guide](Command_Line_Tools.md)** - Complete CLI documentation
- **[Duplicate Cleanup Guide](Duplicate_Cleanup_Guide.md)** - Detailed duplicate management guide

### API Documentation
- **[External API Integration](External_API.md)** - External service integrations
- **[Tenant API](tenant_api.md)** - Tenant management API reference
- **[User Management Integration](user_management_integration.md)** - User service integration

### Architecture & Design
- **[Tenant & User Management](tenant_user.md)** - Multi-tenant architecture overview

## 🚀 Quick Start

### For Developers
1. Read the [main README](../README.md) for project setup
2. Use [CLI Quick Reference](CLI_Quick_Reference.md) for common tasks
3. Check [Command Line Tools Guide](Command_Line_Tools.md) for detailed usage

### For System Administrators
1. Review [Command Line Tools Guide](Command_Line_Tools.md) for maintenance operations
2. Study [Duplicate Cleanup Guide](Duplicate_Cleanup_Guide.md) for data management
3. Understand [External API Integration](External_API.md) for service dependencies

### For DevOps Engineers
1. Check [CLI Quick Reference](CLI_Quick_Reference.md) for automation scripts
2. Review [Command Line Tools Guide](Command_Line_Tools.md) for CI/CD integration
3. Study database management procedures

## 📚 Documentation Categories

### 🛠️ Development Tools
| Document | Purpose | Audience |
|----------|---------|----------|
| [CLI Quick Reference](CLI_Quick_Reference.md) | Essential commands | All users |
| [Command Line Tools Guide](Command_Line_Tools.md) | Complete CLI documentation | Developers, DevOps |
| [Duplicate Cleanup Guide](Duplicate_Cleanup_Guide.md) | Data management | Administrators |

### 🔗 API & Integration
| Document | Purpose | Audience |
|----------|---------|----------|
| [External API](External_API.md) | Service integrations | Developers |
| [Tenant API](tenant_api.md) | Tenant management | Backend developers |
| [User Management Integration](user_management_integration.md) | User service | Backend developers |

### 🏗️ Architecture
| Document | Purpose | Audience |
|----------|---------|----------|
| [Tenant & User Management](tenant_user.md) | Multi-tenant design | Architects, Senior developers |

## 🎯 Common Use Cases

### Database Management
```bash
# Setup new environment
go run helpers/cmd/cmd.go -action=migrate
go run helpers/cmd/cmd.go -action=seed

# Regular maintenance
go run helpers/cmd/cmd.go -action=cleanup-duplicates -dry-run
```
**📖 See**: [Command Line Tools Guide](Command_Line_Tools.md)

### Duplicate Document Cleanup
```bash
# Check for duplicates
go run helpers/cmd/cmd.go -action=cleanup-duplicates -dry-run

# Clean specific asset
go run helpers/cmd/cmd.go -action=cleanup-duplicates -asset-id=UUID
```
**📖 See**: [Duplicate Cleanup Guide](Duplicate_Cleanup_Guide.md)

### API Integration
Understanding external service dependencies and authentication flows.

**📖 See**: [External API Integration](External_API.md)

## 🔍 Finding Information

### By Role

**Developers**:
- Start with [README](../README.md)
- Essential commands: [CLI Quick Reference](CLI_Quick_Reference.md)  
- API integration: [External API](External_API.md)

**System Administrators**:
- Database management: [Command Line Tools Guide](Command_Line_Tools.md)
- Data cleanup: [Duplicate Cleanup Guide](Duplicate_Cleanup_Guide.md)
- Service monitoring: [External API](External_API.md)

**DevOps Engineers**:
- Automation: [CLI Quick Reference](CLI_Quick_Reference.md)
- CI/CD integration: [Command Line Tools Guide](Command_Line_Tools.md)
- Service dependencies: [External API](External_API.md)

### By Task

| Task | Documentation |
|------|---------------|
| Setup new environment | [README](../README.md) + [Command Line Tools](Command_Line_Tools.md) |
| Clean duplicate documents | [Duplicate Cleanup Guide](Duplicate_Cleanup_Guide.md) |
| Database maintenance | [Command Line Tools Guide](Command_Line_Tools.md) |
| API integration | [External API](External_API.md) |
| Troubleshooting | All guides have troubleshooting sections |

## 📝 Contributing to Documentation

When adding new documentation:

1. **Update this index** with new documents
2. **Cross-reference** related documents
3. **Include examples** and use cases
4. **Add troubleshooting** sections
5. **Test all commands** and examples

## 📞 Support

If you can't find what you're looking for:

1. Check the troubleshooting sections in relevant guides
2. Review error messages and logs
3. Verify environment configuration
4. Consult the main [README](../README.md)

---

**Last Updated**: May 28, 2025  
**Version**: 1.0  
**Maintained by**: LecSens Development Team
