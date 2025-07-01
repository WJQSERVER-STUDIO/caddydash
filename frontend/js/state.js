// js/state.js - 管理应用的共享状态

export const state = {
    isEditing: false,
    initialFormState: '', // 用于检测表单是否有未保存的更改
    availableTemplates: [], // 存储从后端获取的可用模板名称
    headerPresets: [],
};