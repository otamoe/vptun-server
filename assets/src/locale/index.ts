import { App, inject, unref } from 'vue'


export interface Language {
    level: number
    name: string
    code: string
    codes: ({code: string})[]
}

export interface Translation { [path: string]: { [value: string]: string }}

export declare interface LocaleInterface {
    translation: (path: string, value: string, defaultValue?: string)=> string;
    setLanguage:(code: string) => boolean;
    load: () => Promise<void>;
    install:(app: App) => void;
}



declare module '@vue/runtime-core' {
  export interface ComponentCustomProperties {
    $t: (path: string, value: string, defaultValue?: string) => string;
    $locale: LocaleInterface;
  }
} 

export class Locale  {
    public languages: Language[] = [
        { 
            level: 0,
            name: 'English',
            code: 'en',
            codes: [
                { code: 'en-Latn' },
                { code: 'en-Latn-US' },
                { code: 'en-Latn-GB' },
                { code: 'en-US' },
                { code: 'en-GB' },
                { code: 'en' },
            ],
        },
        {
            level: 0,
            name: '中文（繁體）',
            code: 'zh-Hant',
            codes: [
                { code:'zh-Hant-HK'},
                { code:'zh-Hant-TW'},
                { code:'zh-Hant-MO'},
                { code:'zh-Hant'},
                { code:'zh-HK'},
                { code:'zh-TW'},
                { code:'zh-MO'},
            ],
        },
        {
            level: 0,
            name: '中文（简体）',
            code: 'zh-Hans',
            codes: [
                { code:'zh-Hans-CN' },
                { code:'zh-Hans-SG' },
                { code:'zh-Hans' },
                { code:'zh-CN' },
                { code:'zh-SG' },
                { code:'zh' },
            ],
        },
        {
            level: 0,
            name: '日本語',
            code: 'ja',
            codes: [
                { code: 'ja-Jpan-JP' },
                { code: 'ja-Jpan' },
                { code: 'ja' },
            ],
        },
    ]

    public language: Language = this.languages[0]
    
    public translations: { [code: string]: Translation} = {}

    private installed = false

    public constructor(code?: string) {
        if (code) {
            this.setLanguage(code)
        }
    }
    
    public translation(path: string, value: string, params?: {[key:string]: any}): string {
        if (!this.translations[this.language.code]) {
            return value
        }
        let result: string
        if (!this.translations[this.language.code][path]) {
            result = value
        } else {
            result = this.translations[this.language.code][path][value]  || value
        }
        if (params) {
            const keys = []
            for (const key in params) {
                keys.push(key.replace(/[.*+?^=!:${}()|[\]\/\\]/g, '\\$&'))
            }
            
            /* eslint prefer-regex-literals: off */
            const re = new RegExp('('+ keys.join("|") +')', 'g')
            result = result.replace(re, function(substring: string, ...args: any[]): string {
                return params[substring] || ''
            })
        }
        return result
    }

    public setLanguage(code: string): boolean {
        const languages:{ [code: string]: Language} = {}
        for (let i = 0; i < this.languages.length; i++) {
            const language = this.languages[i];
            for (let ii = 0; ii < language.codes.length; ii++) {
                const code = language.codes[ii];
                languages[code.code.toLowerCase()] = language
            }
        }
        const key: string = code.replace(/_/g, '-').toLowerCase()
        if (!languages[key]) {
            return false
        }
        this.language = languages[key]
        return true
    }


    public async load(): Promise<void> {
        let translation: Translation = {}
        switch (this.language.code) {
            case 'en':
                translation = (await import('@/locale/en')).default
                break;
            case 'ja':
                translation = (await import('@/locale/ja')).default
                break;
            case 'zh-Hant':
                translation = (await import('@/locale/zh-Hant')).default
                break;
            case 'zh-Hans':
                translation = (await import('@/locale/zh-Hans')).default
                break;
        }
        this.translations[this.language.code] = translation
    }
    public install(app: App): void {
        if (this.installed) {
            return
        }
        this.installed = true
        install(this, app)
    }
}


function install(locale: Locale, app: App): void {
    app.config.globalProperties.$t = function(value : string, params?: {[key:string]: any}): string {
        return locale.translation(this.$.type.name, value, params)
    }
    app.config.globalProperties.$translation = function(path: string, value : string, params?: {[key:string]: any}): string {
        return locale.translation(path, value, params)
    }
    
    app.config.globalProperties.$locale = locale
    app.provide("locale", locale)
    app.provide("t", app.config.globalProperties.$t)
}


export function createLocale(): Locale {
    return new Locale()
}




const defaultLocale = new Locale()
export function useLocale(): Locale {
  return inject("locale", defaultLocale)
}
export function useTranslation(): (path: string, value : string, params?: {[key:string]: any}) => string {
  return inject("translation", defaultLocale.translation)
}
