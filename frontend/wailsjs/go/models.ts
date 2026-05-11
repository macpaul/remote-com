export namespace engine {
	
	export class SerialConfig {
	    baudRate: number;
	    dataBits: number;
	    parity: string;
	    stopBits: string;
	    flowControl: string;
	    charDelay: number;
	    lineDelay: number;
	
	    static createFrom(source: any = {}) {
	        return new SerialConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.baudRate = source["baudRate"];
	        this.dataBits = source["dataBits"];
	        this.parity = source["parity"];
	        this.stopBits = source["stopBits"];
	        this.flowControl = source["flowControl"];
	        this.charDelay = source["charDelay"];
	        this.lineDelay = source["lineDelay"];
	    }
	}
	export class SerialPortInfo {
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new SerialPortInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	    }
	}

}

