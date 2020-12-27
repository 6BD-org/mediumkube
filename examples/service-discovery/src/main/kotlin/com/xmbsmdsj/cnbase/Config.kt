package com.xmbsmdsj.cnbase

import org.springframework.context.EnvironmentAware
import org.springframework.core.env.Environment
import org.springframework.core.env.get
import org.springframework.stereotype.Component

@Component
class Config : EnvironmentAware {

    private lateinit var env: Environment

    override fun setEnvironment(env: Environment) {
        this.env = env
    }

    /**
     * Get a list of upstream services
     */
    fun upstreams(): List<String> {
        val dependencyList = env.getProperty("DEPENDENCIES")
        return dependencyList?.split(",")?: emptyList()
    }

    fun myName(): String {
        return env.getProperty("APP_NAME")?:""
    }


}