import { create } from "zustand"
import type { Skill, SkillDetail, SkillStatus } from "@/types"
import * as skillApi from "@/api/skill"

interface SkillState {
  skills: Skill[]
  currentSkill: SkillDetail | null
  skillStatus: SkillStatus | null
  isLoading: boolean
  isUploading: boolean
  uploadProgress: number
  loadSkills: () => Promise<void>
  loadSkill: (id: string) => Promise<void>
  uploadSkill: (name: string, file: File) => Promise<string>
  deleteSkill: (id: string) => Promise<void>
  downloadSkill: (id: string, filename: string) => Promise<void>
  subscribeStatus: (id: string) => () => void
  clearCurrentSkill: () => void
}

export const useSkillStore = create<SkillState>((set) => ({
  skills: [],
  currentSkill: null,
  skillStatus: null,
  isLoading: false,
  isUploading: false,
  uploadProgress: 0,

  loadSkills: async () => {
    set({ isLoading: true })
    try {
      const skills = await skillApi.listSkills()
      set({ skills })
    } catch (error) {
      console.error("Failed to load skills:", error)
    } finally {
      set({ isLoading: false })
    }
  },

  loadSkill: async (id: string) => {
    set({ isLoading: true })
    try {
      const skill = await skillApi.getSkill(id)
      set({ currentSkill: skill })
    } catch (error) {
      console.error("Failed to load skill:", error)
    } finally {
      set({ isLoading: false })
    }
  },

  uploadSkill: async (name: string, file: File) => {
    set({ isUploading: true, uploadProgress: 0 })
    try {
      const result = await skillApi.uploadSkill(name, file, (progress) => {
        set({ uploadProgress: progress })
      })
      return result.skill_id
    } finally {
      set({ isUploading: false, uploadProgress: 0 })
    }
  },

  deleteSkill: async (id: string) => {
    await skillApi.deleteSkill(id)
    set((state) => ({
      skills: state.skills.filter((s) => s.id !== id),
      currentSkill: state.currentSkill?.id === id ? null : state.currentSkill,
    }))
  },

  downloadSkill: async (id: string, filename: string) => {
    await skillApi.downloadSkill(id, filename)
  },

  subscribeStatus: (id: string) => {
    const unsubscribe = skillApi.subscribeSkillStatus(id, (stage, name) => {
      set({ skillStatus: { stage, name } })
    })
    return unsubscribe
  },

  clearCurrentSkill: () => {
    set({ currentSkill: null, skillStatus: null })
  },
}))
