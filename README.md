# chrononoteai
Note taking app for personal use now built in Go. Will integrate with ChatGPT to be able to ask questions about my own notes

## Configuration for the Markdown Directory
- Need to pass in or configure the directory where the markdown files are stored, either via a config file, environment variable, or command-line argument.

## Reading and Parsing the chrononoteai.md Buffer File
- This involves reading the file and splitting notes based on the YAML front matter, which will act as the delimiter.

## Appending to the Correct Markdown Files
- Parse the date from the YAML metadata to determine the appropriate markdown file (e.g., /notes/2024/09/12.md).
- Create the file if it doesn’t already exist and append the note.
- Clearing or Resetting the chrononoteai.md Buffer:
- After successfully processing the notes, you may want to clear the buffer file or move its content to an archive file for future reference.

Example of potential note:

---
title: Meeting with Project Team
date: 2024-09-12
tags:
  - work
  - meeting
  - project
---

Had a productive discussion on the upcoming project milestones. Key points:
- Finalized the scope of the next sprint.
- Discussed potential blockers related to API development.
- Agreed to follow up with the design team by next week.

#Action items:
- Set up meeting with design team.
- Review API documentation by Friday.
