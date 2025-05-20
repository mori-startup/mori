import ollama
# from langchain_ollama import ChatOllama


historique = [{"role": "system", "content": "Vous êtes un assistant utile."}]

def treating_user_request(output:dict):
    client = ollama.Client(host="http://localhost:11434")

    historique.append({"role": "user", "content": output["message"]})

    response = client.chat(
    model="llama3.1:8b",
    messages=historique,
    stream=True,
    )
    stream = ""
    print("response: ", response)
    for chunk in response:
        print("\nchunk: \n",chunk)

        c = chunk["message"]["content"]
        # print(c, end='', flush=True)
        stream += c

        yield chunk
    # return response

    historique.append({"role": "assistant", "content": stream})


def conversation_resumed(output:dict):

    client = ollama.Client(host="http://localhost:11434")


    historique.append({"role": "assistant", "content": output["message"]})
    historique.append({"role": "user", "content": "résume moi ce que nous avons dit précédemment une phrase de maximum 5 mots dans la langue de la conversation."})

    response = client.chat(
    model="llama3.1:8b",
    messages=historique,
    )

    response = response.message.content


async def send_to_openai(message, model="gpt-3.5-turbo"):
    """
    Send a request to OpenAI API and return the response
    
    Args:
        message (str): The user message to send to OpenAI
        model (str): The OpenAI model to use
        
    Returns:
        Generator yielding response chunks
    """
    api_key = os.getenv("OPENAI_API_KEY")
    if not api_key:
        raise HTTPException(status_code=500, detail="OpenAI API key not found")
    
    headers = {
        "Authorization": f"Bearer {api_key}",
        "Content-Type": "application/json"
    }
    
    data = {
        "model": model,
        "messages": [{"role": "user", "content": message}],
        "stream": True
    }
    
    try:
        response = requests.post(
            "https://api.openai.com/v1/chat/completions",
            headers=headers,
            json=data,
            stream=True
        )
        
        if response.status_code != 200:
            error_msg = f"OpenAI API error: {response.status_code} - {response.text}"
            yield {"message": {"content": error_msg}, "created_at": datetime.now(timezone.utc).isoformat()}
            return
            
        # Process the streaming response
        for line in response.iter_lines():
            if line:
                line = line.decode('utf-8')
                if line.startswith('data: '):
                    if line.startswith('data: [DONE]'):
                        break
                    
                    try:
                        json_str = line[6:]  # Remove 'data: ' prefix
                        data = json.loads(json_str)
                        
                        if 'choices' in data and len(data['choices']) > 0:
                            delta = data['choices'][0].get('delta', {})
                            if 'content' in delta and delta['content']:
                                current_time = datetime.now(timezone.utc)
                                yield {
                                    "message": {"content": delta['content']},
                                    "created_at": current_time.isoformat(timespec='seconds').replace('+00:00', 'Z')
                                }
                    except json.JSONDecodeError:
                        continue
                    
    except Exception as e:
        error_msg = f"Error connecting to OpenAI API: {str(e)}"
        current_time = datetime.now(timezone.utc)
        yield {
            "message": {"content": error_msg},
            "created_at": current_time.isoformat(timespec='seconds').replace('+00:00', 'Z')
        }

async def generate_stream(entry_data):
    if entry_data["message"] != "":
        # You can choose which function to use based on a parameter or environment variable
        use_openai = os.getenv("USE_OPENAI", "false").lower() == "true"
        
        if use_openai:
            # Use OpenAI API
            async for chunk in send_to_openai(entry_data["message"]):
                if chunk is not None:
                    output = {
                        "status": "success",
                        "user_id": entry_data["user_id"],
                        "conversation_id": entry_data["conversation_id"],
                        "response": chunk["message"]["content"],
                        "timestamp": chunk["created_at"]
                    }
                    yield f"data: {json.dumps(output)}\n\n"
        else:
            # Use existing treating_user_request function
            for chunk in treating_user_request(entry_data):
                if chunk is not None:
                    output = {
                        "status": "success",
                        "user_id":entry_data["user_id"],
                        "conversation_id": entry_data["conversation_id"],
                        "response": chunk["message"]["content"],
                        "timestamp": chunk["created_at"]
                    }
                    yield f"data: {json.dumps(output)}\n\n"
    else:
        current_time = datetime.now(timezone.utc)
        output = {
            "status": "error",
            "user_id":entry_data["user_id"],
            "conversation_id": entry_data["conversation_id"],
            "response": "user message empty",
            "timestamp": current_time.isoformat(timespec='seconds').replace('+00:00', 'Z')
        }
        yield f"data: {json.dumps(output)}\n\n"