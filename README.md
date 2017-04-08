This is the repository for https://micromdm.io

The website is built with [hugo](https://gohugo.io/). You're welcome to add to it.

# Adding a new blog post

1. Install [hugo](https://gohugo.io/)   
2. Fork this repo and create a new branch (`git checkout -b my_post`)  
3. Create a new blog post: `hugo new blog/my-post-title.md`    
4. Open content/blog/my-post-title.md in your text editor and write an awesome blog post. You can use Markdown syntax for formatting.  
5. Use `hugo server --watch --buildDrafts=true --theme=micromdm` to view and edit this repo. Your local copy of the website will be visible at `http://localhost:1313`  
6. Commit your changes (`git commit -m 'added my-post-title'`) and `git push` your branch.
7. Open a Pull Request.


# Making changes to the site

The templates for the site are all located in the `themes/micromdm` folder. You can view your changes instantly by running
`hugo server --watch --buildDrafts=true --theme=micromdm` and opening `http://localhost:1313` in your browser.

