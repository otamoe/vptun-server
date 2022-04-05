import { App, VNode } from 'vue'

declare global {
    namespace JSX {
        interface Element extends VNode { }
        interface ElementClass extends App { }
        interface IntrinsicElements {
            [elem: string]: any
        }
    }
}
