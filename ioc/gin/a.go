/*
A MVC binder for free-style of function handler with *gin.Context

Abstract

There are may tedious processes for coding on web service:

    1. Type conversion from HTTP query parameter to desired type in GoLang.
    2. Binding body of HTTP POST to JSON object in GoLang.
        2.1. Perform post-process(e.x. trim text) of binding data
        2.2. Perform data validation of binding data
    3. Convert the result data to JSON response.

Gin has provided foundation features on simple and versatile web application,
this framework try to enhance the building of web application on instinct way.
*/
package gin
