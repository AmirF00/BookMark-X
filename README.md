# BookMark-X
/Users/amirrezafadaeizadehbidari/Desktop/github-project/Twits

# Security Tweet Manager

A web application to manage and summarize security-related tweets. Built with Go and vanilla JavaScript.

## Project Structure
```
.
├── README.md
├── main.go
├── static/
│   ├── tweets.json      # Source tweets data
│   └── summaries.json   # Stored summaries and troll marks
└── templates/
    ├── tweets.html      # Tweet listing page
    ├── summary.html     # Summaries overview
    ├── read.html        # Read summaries
    └── tips.html        # Security tips view
```

## Features

### Routes
- **/twitts**: Main page displaying unprocessed tweets
  - Shows statistics (total, trolls, summarized, remaining)
  - Each tweet can be marked as troll or summarized
  - Tweets disappear after being marked or summarized
  - Dark mode UI for better readability

- **/summary**: View all summaries
  - Shows both draft and published summaries
  - Links to original tweets
  - Summary status (draft/published)

- **/read**: Read-only view of summaries
  - Clean view for reading summaries
  - Links to original tweets

- **/tips**: Security tips compilation
  - Shows only the summary content
  - Focused on security-related information

### API Endpoints
- **POST /summary**: Add new summary
  ```json
  {
    "SNum": 1,
    "Handle": "username",
    "Link": "tweet_url",
    "Summary": "summary text",
    "Draft": true
  }
  ```

- **POST /mark-troll**: Mark tweet as troll
  ```json
  {
    "SNum": 1,
    "Handle": "username",
    "Link": "tweet_url"
  }
  ```

## Data Storage
- **tweets.json**: Source tweets data
  ```json
  [{
    "SNum": 1,
    "Handle": "username",
    "Name": "Display Name",
    "ProfilePic": "pic_url",
    "TweetText": "content",
    "TweetLink": "url"
  }]
  ```

- **summaries.json**: Stored summaries and troll marks
  ```json
  [{
    "SNum": 1,
    "Handle": "username",
    "Link": "url",
    "Summary": "content",
    "Draft": true,
    "Troll": false
  }]
  ```

## Setup and Running
1. Clone the repository
2. Place your tweets data in `static/tweets.json`
3. Create an empty `static/summaries.json` file
4. Run the application:
   ```bash
   go run main.go
   ```
5. Access the application at `http://localhost:8080/twitts`

## Dependencies
- Go 1.x
- Standard library only (no external dependencies)

## Features
- Dark mode UI
- Real-time statistics
- In-memory synchronization with file persistence
- Concurrent request handling with mutex protection
- Modal dialog for summary input
- Automatic page refresh after actions
