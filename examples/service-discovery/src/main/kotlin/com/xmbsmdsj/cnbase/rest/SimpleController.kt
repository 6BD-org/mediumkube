package com.xmbsmdsj.cnbase.rest


import com.jakewharton.retrofit2.converter.kotlinx.serialization.asConverterFactory
import com.xmbsmdsj.cnbase.Config
import com.xmbsmdsj.cnbase.client.CNBaseService
import com.xmbsmdsj.cnbase.model.WhoAmI
import kotlinx.serialization.json.Json
import okhttp3.MediaType
import org.slf4j.Logger
import org.slf4j.LoggerFactory
import org.springframework.beans.factory.InitializingBean
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RestController
import retrofit2.Retrofit
import java.lang.Exception
import kotlin.math.log

val logger : Logger = LoggerFactory.getLogger("SimpleController")

@RestController
class SimpleController: InitializingBean {

    @Autowired
    lateinit var config: Config

    val upstreams: MutableList<CNBaseService> = ArrayList()
    lateinit var name: String

    /**
     *
     */
    @GetMapping(path = ["/hello"])
    fun hello(): List<WhoAmI> {
        val res: MutableList<WhoAmI> = ArrayList()
        res.add(
                WhoAmI(name)
        )
        try {
            for (u in upstreams) {
                res.addAll(u.hello().execute().body()?: emptyList())
            }
        } catch (e : Exception) {
            logger.warn("Dependency unavailable")
            return emptyList()
        }

        return res
    }

    override fun afterPropertiesSet() {
        val contentType = MediaType.get("application/json")
        config.upstreams().map {
           val r: Retrofit = Retrofit.Builder()
                   .baseUrl(it)
                   .addConverterFactory(Json {ignoreUnknownKeys=true}.asConverterFactory(contentType))
                   .build()
            upstreams.add(
                    r.create(CNBaseService::class.java)
            )
        }
        name = config.myName()
    }
}