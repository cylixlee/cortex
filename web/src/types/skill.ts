export interface Skill {
  id: string
  name: string
  description: string
  status: string
  stage: number
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
  stage: number
  name: string
}
