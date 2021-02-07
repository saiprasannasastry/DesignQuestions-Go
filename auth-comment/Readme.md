
# Prerequisites
The system assumes the following are available on the system. if not please use the below links to get them
* Docker-compose
  * https://www.edureka.co/community/50699/how-do-i-install-docker-compose-on-linux
  * https://docs.docker.com/docker-for-mac/
* postman(easier to maintain sessions, and this is what I used for testing)
  * https://www.postman.com/downloads/
* The ports 8001 and 5431 should be free.(5431 not necessary but 8001 is a must, since our api response includes local host to container communication)
# Auth system
* after a user registers we give a user back jwt with token validaity of 5 minutes. if you wish to increase the time , change this to your preferred value 
https://github.com/saiprasannasastry/DesignQuestions-Go/blob/5e859670fdd3e82a56733d00008179cf90fcdfb9/auth-comment/jwt.go#L62
* only in the last 30 seconds will the refresh api work and give back another token, otherwise we get an error
# comment system design

* A user is registered to the system.  --> typically implies  a registration to the system

* Signin. --> once the user is registers we signin

* After signing in a user creats a posts to interact

* After creating a post, all registered users can see all the posts --> for our scenario we have an assumption that all users can see all posts( this can get sophisticated in a bigger system , if we consider that the users have a relation (friends), but in our scenario this is out of scope)

* After selecting which post a user wants to interact, a registered user will comment on the post --> a user who has created the post can also be the first person to comment(pretty common in every comment based system) 	

* After commenting on a post, a user can reply to his own comment/add reaction to a different comment. This is achieved by 2 different API's one for subcomment and other for reaction.
  * A little breif on how sub comments was designed.
   *  We all know that a comment replies system follows a tree hirearchy. To map that we have kept the comment as tree in the database as well.
   *  To get all the comments on a post we take a request body with postname and created by, parent_path is optional parameter(why we pass this value is defined in next point). 
   * We do a breadth first search to get all the top level comments and also show the reply count, parent_id. According to me this typically is how a comment system be designed because getting everything in one call is very expensive, instead with the given reply we can take the parent_path and make the next query to get the reply's of one particular comment. Hope this makes sense. And this can go to n level depth. Although the problem statement asked only for one level depth, figured why not :sweat_smile:.
  * The parent_path addition might look hapazord for the API's but hey, if we had a frontend and a UI, the end user would  have not even known these parameters existed .
			
* Now if a user decides  to interact on a comment he does that.In our case by providing the parent_path in the postinteraction API
    
* Any user can select to addreaction to a comment.
  * Again done by using parent path. Looks ugly in the api because of uuid, but will be masked if ui was there.
* A commented_user can delete his comment / a post owner can delete the comment(maybe the postowner did not like the comment/reply ). No one else has the access to. When this happens all comments and subcomments gets deleted.
* A postowner has the right to delete his post. When this happens all comments and subcomments on the post are deleted.
  
# Supported API's
1 http://localhost:8001/register
  * Register registers a user to the system and saves the user in the database . The password is stored in encrypted format
  * Register needs to contain a body with user name and password. screenshot below
 <img src="https://github.com/saiprasannasastry/DesignQuestions-Go/blob/master/auth-comment/images/register.png" width="600" height="400" />

2 http://localhost:8001/signin
  * Signin verifies the username , password stored in the database and if it matches we send back at jwt. For the next API's we can just call those api's and if the jwt is valid we will perform the asked task
  * Signin needs to contain a body with user name and password. Screenshot below
 <img src="https://github.com/saiprasannasastry/DesignQuestions-Go/blob/master/auth-comment/images/signin.png" width="600" height="400" />	

3 http://localhost:8001/createposts
  * After signin a user creates a post. Each post is assosciated with a unique id and also for a little simplification, we have made sure 2 posts (postname ) can't be same for 1 usercreate
  * Create post needs to contain a body with postname. Screenshot below
 <img src="https://github.com/saiprasannasastry/DesignQuestions-Go/blob/master/auth-comment/images/createposts.png" width="600" height="400" />
 
4 http://localhost:8001/getposts
  * Returns a list of all available posts to comment on
  * Does not need any data in the body. Screenshot below
<img src="https://github.com/saiprasannasastry/DesignQuestions-Go/blob/master/auth-comment/images/getposts.png" width="600" height="400" />

5 http://localhost:8001/addcomments
  * Just adds first level comments to the posts
  * Add comment needs to contain a body with postname, createby and the comment the user wants to add the required postname. Can also insert a reaction at the sametime, but does not make sense for the user to like his own comment(although he can do by adding `comment_reaction` in the body). Screenshot below
<img src="https://github.com/saiprasannasastry/DesignQuestions-Go/blob/master/auth-comment/images/addcomments.png" width="600" height="400" />

6 http://localhost:8001/getinteraction
  * Return a response with all the top level comments for a particular post, a parent_path and commented user. Why a parent_path is returned is explained above.
  * Getinteraction requires a body with postname and createdby. Parent_path is optional. But the replies of both (with parent_path and without)are show below. Shows level1 reply and level 2 reply
<img src="https://github.com/saiprasannasastry/DesignQuestions-Go/blob/master/auth-comment/images/get%20interaction.png" width="600" height="400" />
<img src="https://github.com/saiprasannasastry/DesignQuestions-Go/blob/master/auth-comment/images/level1getinteraction.png" width="600" height="400" />
<img src="https://github.com/saiprasannasastry/DesignQuestions-Go/blob/master/auth-comment/images/level2getinteraction.png" width="600" height="400" />

7 http://localhost:8001/postinteraction
  * Post interaction posts reply for a particular comment. This is achieved by addinng parent_path in response
  * Post interaction requires postname, createdby, parent_path , comment(reply). Screenshot below. This also can be replied n levels because of the design. Screenshot below. 2 levels are shown
<img src="https://github.com/saiprasannasastry/DesignQuestions-Go/blob/master/auth-comment/images/level1reply.png" width="600" height="400" />
<img src="https://github.com/saiprasannasastry/DesignQuestions-Go/blob/master/auth-comment/images/level2reply.png" width="600" height="400" />

8 http://localhost:8001/addreaction
  * Add reaction adds for a particular comment
  * Addreaction takes in parentpath and reaction in the body. Sreenshot below
<img src="https://github.com/saiprasannasastry/DesignQuestions-Go/blob/master/auth-comment/images/addreaction.png" width="600" height="400" />

9 http://localhost:8001/deletecomments
  * Delete commets just take the parent_path. If the user trying to delete is not the owner or commented_user. We get a error message. Screenshot below
<img src="https://github.com/saiprasannasastry/DesignQuestions-Go/blob/master/auth-comment/images/failed%20delete%20request.png" width="600" height="400" />
Success scenario
<img src="https://github.com/saiprasannasastry/DesignQuestions-Go/blob/master/auth-comment/images/success%20delete%20comment.png" width="600" height="400" />

10 http://localhost:8001/deleteposts
  * Delete the post and entire comment trace. If the user trying to delete is not the owner, we wont be able to delete the post
<img src="https://github.com/saiprasannasastry/DesignQuestions-Go/blob/master/auth-comment/images/failuredeletepost.png" width="600" height="400" />
<img src="https://github.com/saiprasannasastry/DesignQuestions-Go/blob/master/auth-comment/images/deletepost.png" width="600" height="400" />

11 http://localhost:8001/refresh
  * refreshes the jwt 
