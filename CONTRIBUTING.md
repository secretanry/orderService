# Contributing to WB-L0

Thank you for your interest in contributing to WB-L0! This document provides guidelines and information for contributors.

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL (or use Docker)
- Redis (optional)
- Kafka (or use Docker)

### Development Setup

1. **Fork and Clone**
   ```bash
   git clone https://github.com/your-username/wb-L0.git
   cd wb-L0
   ```

2. **Environment Setup**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Install Dependencies**
   ```bash
   go mod download
   ```

4. **Run Tests**
   ```bash
   go test ./...
   ```

## 🧪 Testing

### Running Tests

```bash
# Unit tests
go test ./...

# Integration tests (requires Docker)
go test -tags=integration ./...

# Monitoring tests
go test ./monitoring_test.go

# Test with coverage
go test -cover ./...
```

### Test Requirements

- Integration tests require Docker containers for PostgreSQL, Redis, and Kafka
- Use `docker-compose -f docker-compose.test.yml up -d` to start test services
- Tests will automatically clean up after completion

## 📝 Code Style

### Go Code

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Run `golint` for code quality checks
- Keep functions small and focused
- Add comments for exported functions and types

### Commit Messages

Use conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test changes
- `chore`: Build/tooling changes

Examples:
```
feat(monitoring): add Prometheus metrics collection
fix(handlers): resolve nil pointer dereference in order retrieval
docs(readme): update installation instructions
```

## 🔧 Development Workflow

### 1. Create Feature Branch

```bash
git checkout -b feature/your-feature-name
```

### 2. Make Changes

- Write code following the style guidelines
- Add tests for new functionality
- Update documentation if needed

### 3. Test Your Changes

```bash
# Run all tests
go test ./...

# Run specific tests
go test ./handlers/...

# Check for race conditions
go test -race ./...

# Run with monitoring
make monitoring-up
go test -tags=integration ./...
```

### 4. Commit and Push

```bash
git add .
git commit -m "feat(scope): your commit message"
git push origin feature/your-feature-name
```

### 5. Create Pull Request

- Use the PR template
- Describe your changes clearly
- Link any related issues
- Ensure all tests pass

## 🐛 Bug Reports

### Before Submitting

1. Check existing issues
2. Try to reproduce the bug
3. Check if it's a configuration issue

### Bug Report Template

```markdown
**Description**
Brief description of the issue

**Steps to Reproduce**
1. Step 1
2. Step 2
3. Step 3

**Expected Behavior**
What should happen

**Actual Behavior**
What actually happens

**Environment**
- OS: [e.g., macOS, Linux, Windows]
- Go version: [e.g., 1.21.0]
- Docker version: [e.g., 20.10.0]

**Additional Information**
Logs, screenshots, etc.
```

## 💡 Feature Requests

### Before Submitting

1. Check if the feature already exists
2. Consider if it fits the project scope
3. Think about implementation complexity

### Feature Request Template

```markdown
**Problem**
Description of the problem this feature would solve

**Proposed Solution**
Description of the proposed solution

**Alternatives Considered**
Other approaches you considered

**Additional Context**
Any other context or screenshots
```

## 🔍 Code Review Process

### Review Checklist

- [ ] Code follows style guidelines
- [ ] Tests are included and passing
- [ ] Documentation is updated
- [ ] No breaking changes (or properly documented)
- [ ] Security considerations addressed
- [ ] Performance impact considered

### Review Guidelines

- Be constructive and respectful
- Focus on the code, not the person
- Suggest improvements, don't just point out issues
- Use inline comments for specific suggestions

## 🚀 Release Process

### Versioning

We use [Semantic Versioning](https://semver.org/):
- `MAJOR.MINOR.PATCH`
- `MAJOR`: Breaking changes
- `MINOR`: New features (backward compatible)
- `PATCH`: Bug fixes (backward compatible)

### Release Checklist

- [ ] All tests passing
- [ ] Documentation updated
- [ ] Changelog updated
- [ ] Version tagged
- [ ] Docker images built and pushed
- [ ] Release notes written

## 📚 Documentation

### Documentation Standards

- Keep documentation up to date
- Use clear, concise language
- Include examples where helpful
- Update README.md for user-facing changes
- Update technical docs for implementation changes

### Documentation Structure

- `README.md`: Project overview and quick start
- `MONITORING.md`: Monitoring setup and usage
- `TESTING.md`: Testing guidelines and examples
- `GRAFANA-SETUP.md`: Grafana configuration
- `CONTRIBUTING.md`: This file

## 🤝 Community Guidelines

### Code of Conduct

- Be respectful and inclusive
- Help others learn and grow
- Focus on constructive feedback
- Respect different perspectives and experiences

### Communication

- Use GitHub issues for bug reports and feature requests
- Use GitHub discussions for questions and ideas
- Be patient with responses
- Help others when you can

## 🎯 Areas for Contribution

### High Priority

- Bug fixes
- Performance improvements
- Security enhancements
- Documentation improvements

### Medium Priority

- New monitoring metrics
- Additional test coverage
- Code refactoring
- CI/CD improvements

### Low Priority

- New features (discuss first)
- Major architectural changes (discuss first)
- Dependency updates

## 📞 Getting Help

- Check existing documentation
- Search existing issues
- Create a new issue for bugs
- Start a discussion for questions
- Join our community channels

Thank you for contributing to WB-L0! 🎉 