export namespace main {
	
	export class SystemInfo {
	    totalRAM: number;
	    is32Bit: boolean;
	    defaultRAM: number;
	    minRAM: number;
	    maxRAM: number;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalRAM = source["totalRAM"];
	        this.is32Bit = source["is32Bit"];
	        this.defaultRAM = source["defaultRAM"];
	        this.minRAM = source["minRAM"];
	        this.maxRAM = source["maxRAM"];
	    }
	}

}

