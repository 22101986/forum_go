# Forum Project

## Overview
This project is a web-based forum application built using **Go** and **SQLite**. The forum enables users to communicate, share ideas, and interact through posts and comments while offering features like user authentication, likes and dislikes, filtering posts, and Docker containerization.

## Features
### Authentication
- User registration with email, username, and password.
- Unique email verification to prevent duplicate registrations.
- Encrypted password storage (Bonus).
- User login sessions using cookies with expiration dates.
- Single active session per user.

### Communication
- Registered users can create posts and comments.
- Posts can be associated with one or more categories.
- All users (registered or not) can view posts and comments.

### Likes and Dislikes
- Registered users can like or dislike posts and comments.
- The total number of likes and dislikes is visible to all users.

### Filtering
- Filter posts by:
  - Categories (Subforums)
  - User's created posts (For registered users only)
  - Liked posts (For registered users only)

### Notifications
-For likes and dislike
-For moderator asking 

###Guests
-These are unregistered-users that can neither post, comment, like or dislike a post. They only have the permission to see those posts, comments, likes or dislikes.

###Users
-These are the users that will be able to create, comment, like or dislike posts.

###Moderators
-Moderators, as explained above, are users that have a granted access to special functions :
They should be able to monitor the content in the forum by deleting or reporting post to the admin
To create a moderator the user should request an admin for that role

###Administrators
-Users that manage the technical details required for running the forum. This user must be able to :
Promote or demote a normal user to, or from a moderator user.
Receive reports from moderators. If the admin receives a report from a moderator, he can respond to that report
Delete posts and comments
Manage the categories, by being able to create and delete them.

## Database
The forum uses **SQLite** to store data, including:
- Users
- Posts
- Comments
- Categories
- Likes/Dislikes

### Queries
The following SQL queries are used:
- `SELECT` to retrieve data.
- `CREATE` to create tables.
- `INSERT` to add data to tables.

## Docker
The application is containerized using **Docker**, ensuring compatibility and simplifying the deployment process.

### Docker Features
- Containerizing the web forum.
- Environment configuration.
- Easy setup and deployment.

## Error Handling
- HTTPS status management.
- User-friendly error messages.
- Technical error detection and management.

## Technologies Used
- **Go** (Standard packages)
- **SQLite**
- **bcrypt** for password encryption (Bonus)
- **gofrs/uuid** or **google/uuid** for session management (Bonus)
- **Docker**

## Acknowledgements
- SQLite Documentation
- Go Documentation
- Docker Documentation

