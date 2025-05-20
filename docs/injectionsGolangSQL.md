# SQL Injection Prevention in Go

SQL injection is a common attack vector where an attacker can inject malicious SQL code into a query. In Go, one of the primary ways to prevent SQL injection is by using parameterized queries or prepared statements. This document explains how our Go code defends against SQL injections.

## Parameterized Queries

Parameterized queries separate the SQL command from the data. Instead of concatenating user input directly into the SQL string, placeholders (like `$1`, `$2`, etc.) are used. The database driver then safely escapes the provided values.

### Example from the Code

Consider the following snippet from our **UserRepository** in `server.go`:

```go
query := `
    INSERT INTO users (
        user_id, email, first_name, last_name, nickname, about, 
        password, birthday, image, verification_token, verified
    )
    VALUES(
        $1, $2, $3, $4, NULLIF($5, ''), $6,
        $7, $8, $9, $10, $11
    );
`
_, errDB := repo.DB.Exec(
    query,
    user.ID,
    user.Email,
    user.FirstName,
    user.LastName,
    user.Nickname,
    user.About,
    user.Password,
    user.DateOfBirth,
    user.ImagePath,
    user.VerificationToken,
    user.Verified,
)
