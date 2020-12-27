package com.xmbsmdsj.cnbase.client

import com.xmbsmdsj.cnbase.model.WhoAmI
import retrofit2.Call
import retrofit2.http.GET

interface CNBaseService {
    @GET("/hello")
    fun hello(): Call<List<WhoAmI>>

}