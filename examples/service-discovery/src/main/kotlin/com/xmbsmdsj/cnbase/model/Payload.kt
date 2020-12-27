package com.xmbsmdsj.cnbase.model

import kotlinx.serialization.Serializable
import java.time.LocalDateTime

@Serializable
data class WhoAmI(
        val name: String
)