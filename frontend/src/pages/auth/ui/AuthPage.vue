<template>
  <AuthShell :title="flow.title" :subtitle="flow.subtitle">
    <AuthPhoneForm v-if="flow.step === 'phone'" />
    <AuthOtpForm v-else-if="flow.step === 'otp'" />
    <AuthPasswordForm v-else-if="flow.step === 'password'" />
    <AuthUsernameForm v-else-if="flow.step === 'username'" />
  </AuthShell>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import AuthShell from '@/features/auth/ui/AuthShell.vue'
import AuthPhoneForm from '@/features/auth/ui/AuthPhoneForm.vue'
import AuthOtpForm from '@/features/auth/ui/AuthOtpForm.vue'
import AuthPasswordForm from '@/features/auth/ui/AuthPasswordForm.vue'
import AuthUsernameForm from '@/features/auth/ui/AuthUsernameForm.vue'
import { useAuthFlowStore } from '@/features/auth/model/auth-flow.store'

const flow = useAuthFlowStore()
onMounted(() => flow.syncStepFromSession())
</script>
