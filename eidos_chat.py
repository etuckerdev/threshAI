import subprocess  
from personality_engine import PersonalityEngine  

def chat_with_eidos():  
    chat_history = []  # Store chat history in memory  
    personality = PersonalityEngine()  
    print("Eidos: Hi Plato! Let’s tackle this together. What’s on your mind?")  

    while True:  
        try:  
            user_input = input("\nYou: ")  
            if user_input.lower() in ["exit", "quit"]:  
                print("Eidos: Catch you later! Let me know if you need anything.")  
                save_chat_history(chat_history)  # Save history before exiting  
                break  

            # Call the thresh generate command  
            command = ["./thresh", "generate", user_input]  
            result = subprocess.run(command, capture_output=True, text=True)  

            if result.returncode == 0:  
                # Extract only the generated text (skip debug logs)  
                generated_text = result.stdout.split("Generated: ")[-1].strip()  

                # Apply personality tweaks  
                response = personality.generate_response(generated_text)  
                print(f"\n{response}")  

                # Add to chat history  
                chat_history.append(f"You: {user_input}")  
                chat_history.append(response)  
            else:  
                print(f"\nEidos: Oops, something went wrong. Error: {result.stderr.strip()}")  
        except KeyboardInterrupt:  
            print("\nEidos: Caught interrupt. Saving chat history and exiting...")  
            save_chat_history(chat_history)  
            break  

def save_chat_history(chat_history):  
    with open("chat_history.txt", "w") as f:  
        f.write("\n".join(chat_history))  
    print("\nChat history saved to chat_history.txt!")  

if __name__ == "__main__":  
    chat_with_eidos()  