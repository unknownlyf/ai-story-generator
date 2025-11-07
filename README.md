🎬 AI Video Story Generator (Corrected and Stabilized Iteration)

Proprietor: Shreyas Satpute

Repository Location: https://github.com/unknownlyf/ai-story-generator

Designated Methodology: Docker Compose utilization supplementary to local compilation

🎯 Objective

The present undertaking constitutes a polyglot, artificially intelligent application engineered for the automated generation of abbreviated video narratives. Generation is predicated upon textual prompts provided by an operator. The architecture leverages a pipeline of containerized microservices to effectuate narrative generation (Ollama), auditory synthesis (Python), and video-graphic assembly (FFmpeg), orchestrated by a backend implemented in Go. This repository contains a fully rectified and stabilized iteration of the progenitor concept, wherein numerous anomalies at the code-level and environmental-level have been remediated.

✅ Enumeration of Functional Capabilities

Artificial Intelligence Narrative Generation: Employment of the orca-mini Large Language Model, facilitated by the Ollama framework, for the composition of original narratives.

Text-to-Speech Synthesis: Automated conversion of the generated textual narrative into a corresponding auditory voice-over.

Subtitle Generation: Utilization of the OpenAI Whisper model for the transcription of the aforementioned auditory stream into temporally synchronized captions. 

Automated Video-graphic Assembly: The amalgamation of a randomized background video asset, the artificially generated audio, and the resultant subtitles into a final .mp4 container.

Containerized Execution Environment: Comprehensive encapsulation utilizing Docker Compose to ensure a deterministic, cross-platform, and reproducible operational state.

## 🏗️ Architectural Layout of the Project Directory

```text
ai-story-generator/
├── .github/
│   └── workflows/
│       └── … (CI/CD not implemented)
├── internal/
│   └── processes/
│       ├── generation.go
│       ├── speech.go
│       └── video.go        # (Rectified Go source code)
├── resources/
│   └── videos/
│       └── 1.mp4           # Background video asset
├── scripts/
│   ├── captions.py         # (Rectified Python source code)
│   ├── requirements.txt
│   └── tts.py
├── Dockerfile.ollama
├── Dockerfile.worker       # (Rectified Dockerfile)
├── docker-compose.yaml
├── go.mod
├── go.sum
├── setup-ollama.sh
└── README.md               # Project documentation


🛠️ Compendium of Utilized Technologies

Backend Orchestration: Go (Golang)

Artificial Intelligence & Machine Learning Modalities: Python, Ollama, OpenAI Whisper, MoviePy

Environment Encapsulation: Docker, Docker Compose

Core Instrumentation: FFmpeg, Git

Base Container Images: golang:1.22-alpine, python:3.12-slim

🚀 Protocol for System Initialization, Compilation, and Execution (Docker)

It is stipulated that this project be executed exclusively utilizing Docker. The local setup methodology is expressly contraindicated owing to inherent environmental complexities.

Prerequisites for Operation

Confirmation of the installation and operational status of the following software components is mandatory:

Git

The Ollama Application (Host-resident)

Procurement: https://ollama.com/download

Docker Desktop (Host-resident)

Procurement: https://www.docker.com/products/docker-desktop/

Installation and Execution Protocol

Repository Acquisition:

git clone [https://github.com/unknownlyf/ai-story-generator.git](https://github.com/unknownlyf/ai-story-generator.git)
cd ai-story-generator


Service Initiation:
Verification of the operational status of both the Ollama and Docker Desktop host applications is required.

Artificial Intelligence Model Procurement:
A one-time execution of the following command is necessary to retrieve the orca-mini model.

ollama run orca-mini


(Continuation is contingent upon download completion. The terminal session may then be terminated.)

Container Compilation and Instantiation:
The following command initiates the application build process utilizing the rectified artifacts, followed by server instantiation. The initial build cycle requires a non-trivial time allocation.

docker compose up --build


A successful instantiation is indicated by the log message:
worker  | {"time":"...","message":"Server running on port 8080"}

📊 Operational Use and Subsequent Verification

Invocation of Generation Request

The terminal session executing docker compose must be maintained in an active state.

A secondary, independent terminal session must be initiated.

The following curl command is to be executed to transmit the generation prompt to the server endpoint:

curl --location 'http://localhost:8080/api/v1/generate' \
--header 'Content-Type: application/json' \
--data '{
    "message": "Tell me a short story about success after many failures."
}'


Anticipated System Response

The initiating curl terminal will return a JSON-formatted response upon successful completion, enumerating the path of the generated artifact:
{"video":"generated/noggvtl.mp4"} (The specific filename is generated randomly).

The docker compose terminal will display comprehensive logs detailing the sequential execution of the pipeline stages.

The final .mp4 artifact shall materialize within the generated directory of the project structure.

🔧 Compendium of Anomaly Remediation Protocols

This project necessitated extensive debugging activities. The subsequent sections detail common anomalies encountered and their respective resolutions, which have been integrated into the present iteration.

Anomaly 1: make: command not found

Condition: The antecedent documentation suggested the use of the make start command.
Etiology: The make utility is not a standard component of the Windows operating system.
✅ Resolution: Direct invocation of the docker compose command is the prescribed alternative.

docker compose up --build


Anomaly 2: Build failure due to $'\r': command not found

Condition: The Docker build process terminates prematurely during the execution of setup-ollama.sh.
Etiology: The script artifact possesses Windows-formatted (CRLF) line endings, which are non-compliant with the Linux container environment.
✅ Resolution: Application of the dos2unix utility is required to reformat the artifact, followed by a subsequent build attempt.

dos2unix ./setup-ollama.sh
docker compose up --build


Anomaly 3: Build failure due to git not found or moviepy not found

Condition: The pip install procedure fails during the Docker image construction phase.
Etiology: The progenitor Dockerfile.worker artifact contained procedural ordering deficiencies. Specifically, git was invoked prior to its installation, and the moviepy dependency was omitted entirely.
✅ Resolution: This anomaly has been preemptively rectified in the current repository iteration. The Dockerfile.worker now enforces the correct sequence of operations: system-level dependencies (git, ffmpeg) are installed prior to the establishment of the Python virtual environment and the subsequent installation of application-level dependencies.

Anomaly 4: exit status 1 from captions.py (Runtime Anomaly)

Condition: Post-submission of a curl request, the server logs indicate error running autocap: exit status 1.
Etiology: This constituted the principal system defect. The Go application (video.go) failed to create the target directory (./generated) prior to instructing the Python script (captions.py) to write artifacts therein.
✅ Resolution: This anomaly has been preemptively rectified in the current repository iteration.

The internal/processes/video.go artifact has been modified to include an os.MkdirAll("./generated", 0755) invocation prior to the script execution call.

The scripts/captions.py artifact has been substantially rewritten to correctly parse the file paths and arguments as supplied by the rectified Go process.

Anomaly 5: Docker Desktop Initiation Failure

Condition: The Docker Desktop application fails to launch, frequently referencing WSL.
Etiology: The Windows Subsystem for Linux (WSL) is either superannuated or in a corrupted state.
✅ Resolution: Execution of wsl --update within an administrative PowerShell session is required, followed by a mandatory system restart.

🎓 Index of Acquired Competencies

Full-Stack Anomaly Diagnosis: The tracing of a singular error (exit status 1) from its API manifestation, through the Go backend, into the Docker container, ultimately identifying a file system I/O error in a Python script.

Containerization Best Practices: The refactoring of a deficient Dockerfile into a stable, multi-stage compilation process that correctly manages both system-level (apt-get) and application-level (pip) dependencies within an isolated virtual environment.

Polyglot System Integration: The management of the interoperational "contract" between a Go application and its subordinate Python scripts, ensuring the fidelity of file paths, permissions, and runtime arguments.

DevOps Problem Resolution: The systematic remediation of critical environmental impediments (WSL, line-endings) and code-level logical defects, facilitating the transformation of a non-functional artifact into a reliable, operational application.

🔗 Supplemental Resources

Source Artifact: https://github.com/unknownlyf/ai-story-generator

📞 Administrative Contact

Proprietor: Shreyas Satpute

Project Designation: AI Story Generator (Rectified Version)

<sub align="right">Derivative work based on the original project by ccallazans</sub>
