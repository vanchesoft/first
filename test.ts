import { z } from 'zod';
import { OpenAI } from 'langchain/llms/openai';
import { PromptTemplate } from 'langchain/prompts';
import {
    StructuredOutputParser,
    OutputFixingParser,
} from 'langchain/output_parsers';

const parser = StructuredOutputParser.fromZodSchema(
    z.object({
        name: z.string().describe('Human name'),
        surname: z.string().describe('Human surname'),
        age: z.number().describe('Human age'),
        appearance: z.string().describe('Human appearance description'),
        shortBio: z.string().describe('Short bio secription'),
        university: z.string().optional().describe('University name if attended'),
        gender: z.string().describe('Gender of the human'),
        interests: z
            .array(z.string())
            .describe('json array of strings human interests'),
    })
);

const formatInstructions = parser.getFormatInstructions();

const prompt = new PromptTemplate({ 
    template: `Generate details of a hypothetical person.\n{format_instructions} Person description: {description}`, 
    inputVariables: ["description"], 
    partialVariables: { format_instructions: formatInstructions }, 
});

const model = new OpenAI({ temperature: 0.5, model: "gpt-3.5-turbo" }, { basePath: "http://localhost:1337/v1" }); 
const input = await prompt.format({ description: "A man, living in Poland", }); 
const response = await model.call(input);

console.log(prompt)
