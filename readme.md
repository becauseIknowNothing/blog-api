**A blog api to do basic CRUD operations**
**USAGE :** 
 1. */*  Returns testing json, can be used to check working of database (GET)
 2. ***/blogs*** Return all the blogs (GET)
 3.  ***/createblog*** Can be used to create blog and add to the databse (POST)
 4. ***/readblog/{title}*** Can be used to fetch all the blogs having value of title parameter as the substring of blogs titles. (GET)
 5. ***/updateblog/{title}*** Can be used to update the blog where title is known. (POST)
 6. ***/deleteblog/{title}*** Can be used to delete the blog where  title is known. (DELETE)
