import React from 'react'
import {posts} from "../assets/posts-data.js"
import {comments} from "../assets/comments-data.js"
import Post from "../components/Post.jsx"

function Home() {
  return (
    <div className='posts'>
        {posts.map((it,idx)=>{
            return <Post key={idx} post={it}/>
        })}
    </div>
  )
}

export default Home