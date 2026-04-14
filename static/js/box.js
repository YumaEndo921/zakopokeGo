async function releasePokemon(id, name) {
    if (!confirm(`本当に ${name} を逃がしますか？\nお別れしちゃうと二度と戻ってこないよ...🥺`)) {
        return;
    }

    const card = document.getElementById(`pokemon-${id}`);
    const messageEl = document.getElementById('farewell-message');

    try {
        const response = await fetch(`/box/release/${id}`, {
            method: 'POST',
        });

        const data = await response.json();

        if (data.status === 'success') {
            // アニメーション開始
            card.classList.add('releasing');

            // メッセージ表示
            setTimeout(() => {
                messageEl.textContent = data.message;
                messageEl.style.display = 'block';
            }, 500);

            // UIから削除
            setTimeout(() => {
                card.remove();
                messageEl.style.display = 'none';

                // ボックスが空になったかチェック
                const remaining = document.querySelectorAll('.pokemon-card').length;
                if (remaining === 0) {
                    location.reload(); // 空の時のメッセージを出すためにリロード
                }
            }, 2500);
        } else {
            alert('エラー: ' + (data.error || '不明なエラーが発生しました'));
        }
    } catch (err) {
        console.error(err);
        alert('通信に失敗しました😭');
    }
}
