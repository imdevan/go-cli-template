const stage = process.env.NODE_ENV || "dev"
const isProduction = stage === "production"
const basePath = isProduction ? "/go-cli-template" : ""

export default {
  url: isProduction ? "https://devan.gg/go-cli-template" : "http://localhost:4321",
  basePath,
  github: "https://github.com/imdevan/go-cli-template/",
  githubDocs: "https://github.com/imdevan/go-cli-template/",
  title: "go-cli-template",
  description: "A generic CLI tool template built with Go, Cobra, and Bubble Tea. This template provides a foundation for building interactive command-line applications with a clean architecture and modern UI components.",
}
