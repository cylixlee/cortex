export interface Skill {
  id: string
  name: string
  description: string
  status: string
  progress: number
  created_at: string
  updated_at: string
}

export interface SkillReference {
  filename: string
  content: string
}

export interface SkillContent {
  overview: string
  references: SkillReference[]
}

export interface SkillDetail extends Skill {
  skill?: SkillContent
}

export interface SkillStatus {
  status: string
  progress: number
}
